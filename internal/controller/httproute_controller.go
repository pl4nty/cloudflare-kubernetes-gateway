package controller

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/dns"
	"github.com/cloudflare/cloudflare-go/v7/zero_trust"
	"github.com/cloudflare/cloudflare-go/v7/zones"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// HTTPRouteReconciler reconciles a HTTPRoute object
type HTTPRouteReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	Namespace string
}

type routeParentGateway struct {
	parentRef *gatewayv1.ParentReference
	gateway   gatewayv1.Gateway
}

// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gatewayclasses,verbs=get
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways,verbs=list
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=httproutes,verbs=get;list;watch
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=httproutes/status,verbs=update
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.2/pkg/reconcile
//
//nolint:gocyclo
func (r *HTTPRouteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// TODO delete DNS records. load all hostnames via tunnel ID in comment? but can't get DNS zone...
	target := &gatewayv1.HTTPRoute{}
	parentGateways := []routeParentGateway{}
	hostnames := []gatewayv1.Hostname{}
	err := r.Get(ctx, req.NamespacedName, target)
	targetFound := err == nil
	if targetFound {
		for _, parentRef := range target.Spec.ParentRefs {
			if err := validateHTTPRouteGatewayParentRef(parentRef); err != nil {
				logger.Info("HTTPRoute parentRef validation failed", "parentRef", parentRef, "error", err)
				result, statusErr := r.updateHTTPRouteAcceptedStatus(ctx, req.NamespacedName, parentRef, metav1.ConditionFalse, gatewayv1.RouteReasonNoMatchingParent, fmt.Sprintf("ParentRef validation failed: %v", err))
				if statusErr != nil || result.Requeue || result.RequeueAfter != 0 {
					return result, statusErr
				}
				return ctrl.Result{}, nil
			}

			namespace := target.Namespace
			if parentRef.Namespace != nil {
				namespace = string(*parentRef.Namespace)
			}
			gateway := &gatewayv1.Gateway{}
			gatewayRef := types.NamespacedName{
				Namespace: namespace,
				Name:      string(parentRef.Name),
			}
			if err := r.Get(ctx, gatewayRef, gateway); err != nil {
				logger.Error(err, "Failed to get Gateway")
				result, statusErr := r.updateHTTPRouteAcceptedStatus(ctx, req.NamespacedName, parentRef, metav1.ConditionFalse, gatewayv1.RouteReasonNoMatchingParent, fmt.Sprintf("No matching parent Gateway %s: %v", gatewayRef.String(), err))
				if statusErr != nil || result.Requeue || result.RequeueAfter != 0 {
					return result, statusErr
				}
				return ctrl.Result{}, err
			}
			parentRef := parentRef
			parentGateways = append(parentGateways, routeParentGateway{parentRef: &parentRef, gateway: *gateway})
		}

		hostnames = target.Spec.Hostnames
	} else {
		if !apierrors.IsNotFound(err) {
			logger.Error(err, "Failed to get HTTPRoute")
			return ctrl.Result{}, err
		}

		gatewayList := &gatewayv1.GatewayList{}
		if err := r.List(ctx, gatewayList); err != nil {
			logger.Error(err, "Failed to list Gateways")
			return ctrl.Result{}, err
		}
		for _, gateway := range gatewayList.Items {
			parentGateways = append(parentGateways, routeParentGateway{gateway: gateway})
		}
	}

	routes := &gatewayv1.HTTPRouteList{}
	if err := r.List(ctx, routes); err != nil {
		logger.Error(err, "Failed to list HTTPRoutes")
		return ctrl.Result{}, err
	}

	for _, parentGateway := range parentGateways {
		gateway := parentGateway.gateway

		// check target is in scope
		gatewayClass := &gatewayv1.GatewayClass{}
		if err := r.Get(ctx, types.NamespacedName{
			Name: string(gateway.Spec.GatewayClassName),
		}, gatewayClass); err != nil {
			logger.Error(err, "Failed to get GatewayClasses")
			return ctrl.Result{}, err
		}

		if gatewayClass.Spec.ControllerName != controllerName {
			continue
		}

		// search for sibling routes
		siblingRoutes := []gatewayv1.HTTPRoute{}
		for _, searchRoute := range routes.Items {
			for _, searchParent := range searchRoute.Spec.ParentRefs {
				namespace := searchRoute.Namespace
				if searchParent.Namespace != nil {
					namespace = string(*searchParent.Namespace)
				}
				if namespace == gateway.Namespace && string(searchParent.Name) == gateway.Name {
					siblingRoutes = append(siblingRoutes, searchRoute)
					break
				}
			}
		}

		// fan out to siblings
		ingress := []zero_trust.TunnelCloudflaredConfigurationUpdateParamsConfigIngress{}
		for _, route := range siblingRoutes {
			for _, rule := range route.Spec.Rules {
				paths := map[string]bool{}
				for _, match := range rule.Matches {
					if match.Path == nil {
						paths["/"] = true
					} else {
						paths[*match.Path.Value] = true
					}

					if match.Headers != nil {
						logger.Info("HTTPRoute header match is not supported", "HTTPRouteMatch.Headers", match.Headers)
					}
				}

				// TODO implement this with rewrite rules? Core filters are a MUST in the spec
				if rule.Filters != nil {
					logger.Info("HTTPRoute filters are not supported", "HTTPRouteFilter", rule.Filters)
				}

				services := map[string]bool{}
				for _, backend := range rule.BackendRefs {
					if backend.Port == nil {
						err := errors.New("HTTPRoute backend port is nil")
						logger.Error(err, "HTTPRoute backend port is required and nil", "backend", backend)
						continue
					}

					var namespace string
					if backend.Namespace == nil {
						namespace = route.Namespace
					} else {
						namespace = string(*backend.Namespace)
					}

					services[fmt.Sprintf("http://%s.%s:%d", backend.Name, namespace, backend.Port)] = true
				}

				// product of hostname, path, service
				for _, hostname := range route.Spec.Hostnames {
					for path := range paths {
						for service := range services {
							ingress = append(ingress, zero_trust.TunnelCloudflaredConfigurationUpdateParamsConfigIngress{
								Hostname: cloudflare.String(string(hostname)),
								Path:     cloudflare.String(path),
								Service:  cloudflare.String(service),
							})
						}
					}
				}
			}
		}

		// last rule must be the catch-all
		ingress = append(ingress, zero_trust.TunnelCloudflaredConfigurationUpdateParamsConfigIngress{
			Service: cloudflare.String("http_status:404"),
		})

		// increment AttachedRoutes in each gateway listener status
		gatewayObj := &gatewayv1.Gateway{}
		gatewayRef := types.NamespacedName{
			Namespace: gateway.Namespace,
			Name:      gateway.Name,
		}
		if err := r.Get(ctx, gatewayRef, gatewayObj); err != nil {
			logger.Error(err, "Failed to re-fetch gateway")
			return ctrl.Result{}, err
		}
		listeners := []gatewayv1.ListenerStatus{}
		for _, listener := range gatewayObj.Status.Listeners {
			listener.AttachedRoutes = int32(len(ingress))
			listeners = append(listeners, listener)
		}
		logger.Info("Updating Gateway listeners", "AttachedRoutes", len(ingress))
		gatewayObj.Status.Listeners = listeners
		if err := r.Status().Update(ctx, gatewayObj); err != nil {
			logger.Error(err, "Failed to update Gateway status")
			return ctrl.Result{}, err
		}

		account, api, err := InitCloudflareAPI(ctx, r.Client, string(gateway.Spec.GatewayClassName))
		if err != nil {
			logger.Error(err, "Failed to initialize Cloudflare API")
			return ctrl.Result{}, err
		}

		tunnels, err := api.ZeroTrust.Tunnels.List(ctx, zero_trust.TunnelListParams{
			AccountID: cloudflare.String(account),
			IsDeleted: cloudflare.Bool(false),
			Name:      cloudflare.String(gateway.Name),
		})
		if err != nil {
			logger.Error(err, "Failed to get tunnel from Cloudflare API")
			return ctrl.Result{}, err
		}
		if len(tunnels.Result) == 0 {
			logger.Info("Tunnel doesn't exist yet, probably waiting for the Gateway controller. Retrying in 1 minute", "gateway", gateway.Name)
			return ctrl.Result{RequeueAfter: time.Minute}, nil
		}
		tunnel := tunnels.Result[0]

		_, err = api.ZeroTrust.Tunnels.Cloudflared.Configurations.Update(ctx, tunnel.ID, zero_trust.TunnelCloudflaredConfigurationUpdateParams{
			AccountID: cloudflare.String(account),
			Config: cloudflare.F(zero_trust.TunnelCloudflaredConfigurationUpdateParamsConfig{
				Ingress: cloudflare.F(ingress),
			},
			),
		})
		if err != nil {
			logger.Error(err, "Failed to update Tunnel configuration")
			return ctrl.Result{}, err
		}

		logger.Info("Updated Tunnel configuration", "ingress", ingress)

		// duplicate CNAMEs can't exist, so the last parentRef wins
		for _, gwHostname := range hostnames {
			hostname := string(gwHostname)
			zoneID, err := FindZoneID(hostname, ctx, api, account)
			if err != nil {
				return ctrl.Result{}, err
			}

			content := fmt.Sprintf("%s.cfargotunnel.com", tunnel.ID)
			comment := "Managed by github.com/pl4nty/cloudflare-kubernetes-gateway"
			records, _ := api.DNS.Records.List(ctx, dns.RecordListParams{
				ZoneID:  cloudflare.String(zoneID),
				Proxied: cloudflare.Bool(true),
				Type:    cloudflare.F[dns.RecordListParamsType]("CNAME"),
				Name:    cloudflare.F(dns.RecordListParamsName{Exact: cloudflare.String(hostname)}),
			})
			if len(records.Result) == 0 {
				_, err := api.DNS.Records.New(ctx, dns.RecordNewParams{
					ZoneID: cloudflare.String(zoneID),
					Body: dns.CNAMERecordParam{
						Proxied: cloudflare.Bool(true),
						Type:    cloudflare.F[dns.CNAMERecordType]("CNAME"),
						Name:    cloudflare.String(hostname),
						Content: cloudflare.String(content),
						Comment: cloudflare.String(comment),
					},
				})
				if err != nil {
					logger.Error(err, "Failed to create DNS record", "hostname", hostname, "content", content)
					return ctrl.Result{}, err
				}
			} else {
				_, err := api.DNS.Records.Update(ctx, records.Result[0].ID, dns.RecordUpdateParams{
					ZoneID: cloudflare.String(zoneID),
					Body: dns.CNAMERecordParam{
						Proxied: cloudflare.Bool(true),
						Type:    cloudflare.F[dns.CNAMERecordType]("CNAME"),
						Name:    cloudflare.String(hostname),
						Content: cloudflare.String(content),
						Comment: cloudflare.String(comment),
					},
				})
				if err != nil {
					logger.Error(err, "Failed to update DNS record", "hostname", hostname, "content", content)
					return ctrl.Result{}, err
				}
			}
		}
		logger.Info("Updated DNS records", "hostnames", hostnames)

		if targetFound && parentGateway.parentRef != nil {
			result, err := r.updateHTTPRouteAcceptedStatus(ctx, req.NamespacedName, *parentGateway.parentRef, metav1.ConditionTrue, gatewayv1.RouteReasonAccepted, "Successfully reconciled with Cloudflare")
			if err != nil || result.Requeue || result.RequeueAfter != 0 {
				return result, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HTTPRouteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	pred := predicate.GenerationChangedPredicate{}
	return ctrl.NewControllerManagedBy(mgr).
		For(&gatewayv1.HTTPRoute{}).
		WithEventFilter(pred).
		Complete(r)
}

func FindZoneID(hostname string, ctx context.Context, api *cloudflare.Client, accountID string) (string, error) {
	logger := log.FromContext(ctx)
	for parts := range len(strings.Split(hostname, ".")) {
		zoneName := strings.Join(strings.Split(hostname, ".")[parts:], ".")
		zoneList, err := api.Zones.List(ctx, zones.ZoneListParams{
			Account: cloudflare.F(zones.ZoneListParamsAccount{ID: cloudflare.String(accountID)}),
			Name:    cloudflare.String(zoneName),
			Status:  cloudflare.F(zones.ZoneListParamsStatusActive),
		})
		if err != nil {
			logger.Error(err, "Failed to list DNS zones")
			return "", err
		}
		if len(zoneList.Result) != 0 {
			return zoneList.Result[0].ID, nil
		}
	}
	err := errors.New("failed to discover DNS zone")
	logger.Error(err, "Failed to discover parent DNS zone. Ensure Zone.DNS permission is configured", "hostname", hostname)
	return "", err
}

func validateHTTPRouteGatewayParentRef(parentRef gatewayv1.ParentReference) error {
	if parentRefGroup(parentRef) != "gateway.networking.k8s.io" {
		return fmt.Errorf("unsupported parentRef group %s", parentRefGroup(parentRef))
	}
	if parentRefKind(parentRef) != "Gateway" {
		return fmt.Errorf("unsupported parentRef kind %s", parentRefKind(parentRef))
	}
	return nil
}

func (r *HTTPRouteReconciler) updateHTTPRouteAcceptedStatus(ctx context.Context, routeRef types.NamespacedName, parentRef gatewayv1.ParentReference, status metav1.ConditionStatus, reason gatewayv1.RouteConditionReason, message string) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	route := &gatewayv1.HTTPRoute{}
	if err := r.Get(ctx, routeRef, route); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("HTTPRoute resource not found. Ignoring status update")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to re-fetch HTTPRoute")
		return ctrl.Result{}, err
	}

	if !setHTTPRouteParentStatusCondition(route, parentRef, metav1.Condition{
		Type:               string(gatewayv1.RouteConditionAccepted),
		Status:             status,
		Reason:             string(reason),
		ObservedGeneration: route.Generation,
		Message:            message,
	}) {
		return ctrl.Result{}, nil
	}

	if err := r.Status().Update(ctx, route); err != nil {
		if apierrors.IsConflict(err) {
			logger.Info("Conflict when updating HTTPRoute status, retrying")
			return ctrl.Result{Requeue: true}, nil
		}
		logger.Error(err, "Failed to update HTTPRoute status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func setHTTPRouteParentStatusCondition(route *gatewayv1.HTTPRoute, parentRef gatewayv1.ParentReference, condition metav1.Condition) bool {
	normalizedParentRef := normalizeHTTPRouteParentRef(parentRef, route.Namespace)

	for i := range route.Status.Parents {
		parentStatus := &route.Status.Parents[i]
		if parentStatus.ControllerName == gatewayv1.GatewayController(controllerName) && parentRefsEqual(normalizeHTTPRouteParentRef(parentStatus.ParentRef, route.Namespace), normalizedParentRef) {
			parentRefChanged := !parentRefsExactlyEqual(parentStatus.ParentRef, normalizedParentRef)
			parentStatus.ParentRef = normalizedParentRef
			if !parentRefChanged && statusConditionMatches(meta.FindStatusCondition(parentStatus.Conditions, condition.Type), condition) {
				return false
			}
			meta.SetStatusCondition(&parentStatus.Conditions, condition)
			return true
		}
	}

	parentStatus := gatewayv1.RouteParentStatus{
		ParentRef:      normalizedParentRef,
		ControllerName: gatewayv1.GatewayController(controllerName),
		Conditions:     []metav1.Condition{},
	}
	meta.SetStatusCondition(&parentStatus.Conditions, condition)
	route.Status.Parents = append(route.Status.Parents, parentStatus)
	return true
}

func normalizeHTTPRouteParentRef(parentRef gatewayv1.ParentReference, defaultNamespace string) gatewayv1.ParentReference {
	if parentRef.Group == nil {
		group := gatewayv1.Group("gateway.networking.k8s.io")
		parentRef.Group = &group
	}
	if parentRef.Kind == nil {
		kind := gatewayv1.Kind("Gateway")
		parentRef.Kind = &kind
	}
	if parentRef.Namespace == nil && defaultNamespace != "" {
		namespace := gatewayv1.Namespace(defaultNamespace)
		parentRef.Namespace = &namespace
	}
	return parentRef
}

func statusConditionMatches(existing *metav1.Condition, condition metav1.Condition) bool {
	return existing != nil &&
		existing.Type == condition.Type &&
		existing.Status == condition.Status &&
		existing.Reason == condition.Reason &&
		existing.ObservedGeneration == condition.ObservedGeneration &&
		existing.Message == condition.Message
}

func parentRefsEqual(left, right gatewayv1.ParentReference) bool {
	return parentRefGroup(left) == parentRefGroup(right) &&
		parentRefKind(left) == parentRefKind(right) &&
		namespacePtrEqual(left.Namespace, right.Namespace) &&
		left.Name == right.Name &&
		sectionNamePtrEqual(left.SectionName, right.SectionName) &&
		portNumberPtrEqual(left.Port, right.Port)
}

func parentRefsExactlyEqual(left, right gatewayv1.ParentReference) bool {
	return groupPtrEqual(left.Group, right.Group) &&
		kindPtrEqual(left.Kind, right.Kind) &&
		namespacePtrEqual(left.Namespace, right.Namespace) &&
		left.Name == right.Name &&
		sectionNamePtrEqual(left.SectionName, right.SectionName) &&
		portNumberPtrEqual(left.Port, right.Port)
}

func parentRefGroup(parentRef gatewayv1.ParentReference) string {
	if parentRef.Group == nil {
		return "gateway.networking.k8s.io"
	}
	return string(*parentRef.Group)
}

func parentRefKind(parentRef gatewayv1.ParentReference) string {
	if parentRef.Kind == nil {
		return "Gateway"
	}
	return string(*parentRef.Kind)
}

func groupPtrEqual(left, right *gatewayv1.Group) bool {
	if left == nil || right == nil {
		return left == right
	}
	return *left == *right
}

func kindPtrEqual(left, right *gatewayv1.Kind) bool {
	if left == nil || right == nil {
		return left == right
	}
	return *left == *right
}

func namespacePtrEqual(left, right *gatewayv1.Namespace) bool {
	if left == nil || right == nil {
		return left == right
	}
	return *left == *right
}

func sectionNamePtrEqual(left, right *gatewayv1.SectionName) bool {
	if left == nil || right == nil {
		return left == right
	}
	return *left == *right
}

func portNumberPtrEqual(left, right *gatewayv1.PortNumber) bool {
	if left == nil || right == nil {
		return left == right
	}
	return *left == *right
}
