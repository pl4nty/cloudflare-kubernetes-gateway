package controller

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cloudflare/cloudflare-go/v2"
	"github.com/cloudflare/cloudflare-go/v2/shared"
	"github.com/cloudflare/cloudflare-go/v2/zero_trust"

	appsv1 "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const gatewayClassFinalizer = "cfargotunnel.com/finalizer"
const gatewayFinalizer = "cfargotunnel.com/finalizer"
const controllerName = "github.com/pl4nty/cloudflare-kubernetes-gateway"

// GatewayReconciler reconciles a Gateway object
type GatewayReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	Recorder  record.EventRecorder
	Namespace string
}

// The following markers are used to generate the rules permissions (RBAC) on config/rbac using controller-gen
// when the command <make manifests> is executed.
// To know more about markers see: https://book.kubebuilder.io/reference/markers.html

// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gatewayclasses,verbs=get;update
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gatewayclasses/finalizers,verbs=update
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways,verbs=get;list;update;watch
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways/finalizers,verbs=update
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways/status,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=create;get;list;update;watch
// +kubebuilder:rbac:groups=core,resources=events,verbs=create
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// It is essential for the controller's reconciliation loop to be idempotent. By following the Operator
// pattern you will create Controllers which provide a reconcile function
// responsible for synchronizing resources until the desired state is reached on the cluster.
// Breaking this recommendation goes against the design principles of controller-runtime.
// and may lead to unforeseen consequences such as resources becoming stuck and requiring manual intervention.
// For further info:
// - About Operator Pattern: https://kubernetes.io/docs/concepts/extend-kubernetes/operator/
// - About Controllers: https://kubernetes.io/docs/concepts/architecture/controller/
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.2/pkg/reconcile
//
//nolint:gocyclo
func (r *GatewayReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the Gateway instance
	// The purpose is check if the Custom Resource for the Kind Gateway
	// is applied on the cluster if not we return nil to stop the reconciliation
	gateway := &gatewayv1.Gateway{}
	if err := r.Get(ctx, req.NamespacedName, gateway); err != nil {
		if apierrors.IsNotFound(err) {
			// If the custom resource is not found then it usually means that it was deleted or not created
			// In this way, we will stop the reconciliation
			log.Info("gateway resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get gateway")
		return ctrl.Result{}, err
	}

	// check if parent GatewayClass is ours and update finalizer
	gatewayClass := &gatewayv1.GatewayClass{}
	if err := r.Get(ctx, types.NamespacedName{Name: string(gateway.Spec.GatewayClassName)}, gatewayClass); err != nil {
		if apierrors.IsNotFound(err) {
			// If the custom resource is not found then it usually means that it was deleted or not created
			// In this way, we will stop the reconciliation
			log.Info("gatewayclass resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get gatewayclass")
		return ctrl.Result{}, err
	}
	if gatewayClass.Spec.ControllerName != controllerName {
		log.Info("Ignoring gateway with non-matching GatewayClass")
		return ctrl.Result{}, nil
	}
	if !controllerutil.ContainsFinalizer(gatewayClass, gatewayClassFinalizer) {
		log.Info("Adding Finalizer for GatewayClass")
		if ok := controllerutil.AddFinalizer(gatewayClass, gatewayClassFinalizer); !ok {
			log.Error(nil, "Failed to add finalizer into the GatewayClass")
			return ctrl.Result{Requeue: true}, nil
		}

		if err := r.Update(ctx, gatewayClass); err != nil {
			log.Error(err, "Failed to update GatewayClass to add finalizer")
			return ctrl.Result{}, err
		}
	}

	account, api, err := InitCloudflareApi(ctx, r.Client, gatewayClass.Name)

	// Let's add a finalizer. Then, we can define some operations which should
	// occur before the custom resource is deleted.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/finalizers
	if !controllerutil.ContainsFinalizer(gateway, gatewayFinalizer) {
		log.Info("Adding Finalizer for Gateway")
		if ok := controllerutil.AddFinalizer(gateway, gatewayFinalizer); !ok {
			log.Error(nil, "Failed to add finalizer into the Gateway")
			return ctrl.Result{Requeue: true}, nil
		}

		if err := r.Update(ctx, gateway); err != nil {
			log.Error(err, "Failed to update Gateway to add finalizer")
			return ctrl.Result{}, err
		}
	}

	// Check if the Gateway instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isGatewayMarkedToBeDeleted := gateway.GetDeletionTimestamp() != nil
	if isGatewayMarkedToBeDeleted {
		if controllerutil.ContainsFinalizer(gateway, gatewayFinalizer) {
			log.Info("Performing Finalizer Operations for Gateway before delete CR")

			// Let's add here a status "Downgrade" to reflect that this resource began its process to be terminated.
			meta.SetStatusCondition(&gateway.Status.Conditions, metav1.Condition{Type: string(gatewayv1.GatewayConditionAccepted),
				Status: metav1.ConditionUnknown, Reason: string(gatewayv1.GatewayReasonPending), ObservedGeneration: gateway.Generation,
				Message: fmt.Sprintf("Performing finalizer operations for the Gateway: %s ", gateway.Name)})

			if err := r.Status().Update(ctx, gateway); err != nil {
				log.Error(err, "Failed to update Gateway finalizer status")
				return ctrl.Result{}, err
			}

			// Perform all operations required before removing the finalizer and allow
			// the Kubernetes API to remove the custom resource.
			if err := r.doFinalizerOperationsForGateway(ctx, gatewayClass, gateway, account, api); err != nil {
				log.Error(err, "Failed to complete finalizer operations for Gateway")
				return ctrl.Result{Requeue: true}, nil
			}

			// Re-fetch the gateway Custom Resource before updating the status
			// so that we have the latest state of the resource on the cluster and we will avoid
			// raising the error "the object has been modified, please apply
			// your changes to the latest version and try again" which would re-trigger the reconciliation
			if err := r.Get(ctx, req.NamespacedName, gateway); err != nil {
				log.Error(err, "Failed to re-fetch gateway")
				return ctrl.Result{}, err
			}

			meta.SetStatusCondition(&gateway.Status.Conditions, metav1.Condition{Type: string(gatewayv1.GatewayConditionAccepted),
				Status: metav1.ConditionTrue, Reason: "Finalizing", ObservedGeneration: gateway.Generation,
				Message: fmt.Sprintf("Finalizer operations for custom resource %s name were successfully accomplished", gateway.Name)})

			if err := r.Status().Update(ctx, gateway); err != nil {
				log.Error(err, "Failed to update Gateway finalizer status")
				return ctrl.Result{}, err
			}

			log.Info("Removing Finalizer for Gateway after successfully perform the operations")
			if ok := controllerutil.RemoveFinalizer(gateway, gatewayFinalizer); !ok {
				log.Error(err, "Failed to remove finalizer for Gateway")
				return ctrl.Result{Requeue: true}, nil
			}

			if err := r.Update(ctx, gateway); err != nil {
				log.Error(err, "Failed to remove finalizer for Gateway")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Validate Gateway listeners and update status
	listenerStatuses := []gatewayv1.ListenerStatus{}
	validListener := false
	for _, listener := range gateway.Spec.Listeners {
		listenerStatus := gatewayv1.ListenerStatus{
			Name:           listener.Name,
			AttachedRoutes: 0,
			SupportedKinds: []gatewayv1.RouteGroupKind{},
		}

		if (listener.Protocol == gatewayv1.HTTPProtocolType && listener.Port == gatewayv1.PortNumber(80)) || (listener.Protocol == gatewayv1.HTTPSProtocolType && listener.Port == gatewayv1.PortNumber(443)) {
			validListener = true
			if listener.TLS != nil && listener.TLS.CertificateRefs != nil {
				ref := listener.TLS.CertificateRefs[0]
				secretRef := types.NamespacedName{
					Name: string(ref.Name),
				}
				if ref.Namespace != nil {
					secretRef.Namespace = string(*ref.Namespace)
				}
				secret := &core.Secret{}

				if err := r.Get(ctx, secretRef, secret); err != nil {
					log.Error(err, "unable to fetch Secret from listener CertificateRefs", "listener", listener.Name)

					meta.SetStatusCondition(&listenerStatus.Conditions, metav1.Condition{
						Type:               string(gatewayv1.ListenerConditionResolvedRefs),
						Status:             metav1.ConditionFalse,
						Reason:             string(gatewayv1.ListenerReasonInvalidCertificateRef),
						ObservedGeneration: gateway.Generation,
						Message:            "Listener TLS certificate references are not supported",
					})
				} else {
					meta.SetStatusCondition(&listenerStatus.Conditions, metav1.Condition{
						Type:               string(gatewayv1.ListenerConditionResolvedRefs),
						Status:             metav1.ConditionFalse,
						Reason:             string(gatewayv1.ListenerReasonRefNotPermitted),
						ObservedGeneration: gateway.Generation,
						Message:            "Listener TLS certificate references are not supported",
					})
				}
			} else {
				meta.SetStatusCondition(&listenerStatus.Conditions, metav1.Condition{
					Type:               string(gatewayv1.ListenerConditionAccepted),
					Status:             metav1.ConditionTrue,
					Reason:             string(gatewayv1.ListenerReasonAccepted),
					ObservedGeneration: gateway.Generation,
					Message:            fmt.Sprintf("Listener protocol %s and port %d accepted", listener.Protocol, listener.Port),
				})
			}
		} else if listener.Protocol != gatewayv1.HTTPProtocolType && listener.Protocol != gatewayv1.HTTPSProtocolType {
			meta.SetStatusCondition(&listenerStatus.Conditions, metav1.Condition{
				Type:               string(gatewayv1.ListenerConditionAccepted),
				Status:             metav1.ConditionFalse,
				Reason:             string(gatewayv1.ListenerReasonUnsupportedProtocol),
				ObservedGeneration: gateway.Generation,
				Message:            fmt.Sprintf("Listener protocol %s is not supported. Only HTTP or HTTPS is supported", listener.Protocol),
			})
		} else if listener.Port != gatewayv1.PortNumber(80) && listener.Port != gatewayv1.PortNumber(443) {
			meta.SetStatusCondition(&listenerStatus.Conditions, metav1.Condition{
				Type:               string(gatewayv1.ListenerConditionAccepted),
				Status:             metav1.ConditionFalse,
				Reason:             string(gatewayv1.ListenerReasonPortUnavailable),
				ObservedGeneration: gateway.Generation,
				Message:            fmt.Sprintf("Listener port %d is not supported. Only port 80 or 443 is supported", listener.Port),
			})
		} else {
			meta.SetStatusCondition(&listenerStatus.Conditions, metav1.Condition{
				Type:               string(gatewayv1.ListenerConditionAccepted),
				Status:             metav1.ConditionFalse,
				Reason:             string(gatewayv1.ListenerReasonInvalid),
				ObservedGeneration: gateway.Generation,
				Message:            "Invalid protocol/port combination. Listener only supports HTTP on port 80 or HTTPS on port 443",
			})
		}

		// search AllowedRoutes for HTTPRoute and default to HTTPRoute if empty
		validKind := false
		for _, kind := range listener.AllowedRoutes.Kinds {
			if kind.Kind == gatewayv1.Kind("HTTPRoute") {
				validKind = true
			} else {
				meta.SetStatusCondition(&listenerStatus.Conditions, metav1.Condition{
					Type:               string(gatewayv1.ListenerConditionResolvedRefs),
					Status:             metav1.ConditionFalse,
					Reason:             string(gatewayv1.ListenerReasonInvalidRouteKinds),
					ObservedGeneration: gateway.Generation,
					Message:            fmt.Sprintf("Listener only supports HTTPRoute, not %s", kind.Kind),
				})
			}
		}
		if validKind || len(listener.AllowedRoutes.Kinds) == 0 {
			listenerStatus.SupportedKinds = []gatewayv1.RouteGroupKind{{Kind: gatewayv1.Kind("HTTPRoute")}}
		}

		listenerStatuses = append(listenerStatuses, listenerStatus)
	}

	if err := r.Get(ctx, req.NamespacedName, gateway); err != nil {
		log.Error(err, "Failed to re-fetch gateway")
		return ctrl.Result{}, err
	}
	gateway.Status.Listeners = listenerStatuses

	if !validListener {
		meta.SetStatusCondition(&gateway.Status.Conditions, metav1.Condition{Type: string(gatewayv1.GatewayConditionAccepted),
			Status: metav1.ConditionFalse, Reason: string(gatewayv1.GatewayReasonListenersNotValid), ObservedGeneration: gateway.Generation,
			Message: fmt.Sprintf("No valid listeners for gateway (%s)", gateway.Name)})

		if err := r.Status().Update(ctx, gateway); err != nil {
			if strings.Contains(err.Error(), "apply your changes to the latest version and try again") {
				log.Info("Conflict when updating Gateway listener status, retrying")
				return ctrl.Result{Requeue: true}, nil
			} else {
				log.Error(err, "Failed to update Gateway listener status")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	meta.SetStatusCondition(&gateway.Status.Conditions, metav1.Condition{Type: string(gatewayv1.GatewayConditionAccepted),
		Status: metav1.ConditionTrue, Reason: string(gatewayv1.GatewayReasonAccepted), ObservedGeneration: gateway.Generation,
		Message: fmt.Sprintf("Validated and accepted gateway (%s)", gateway.Name)})

	if err := r.Status().Update(ctx, gateway); err != nil {
		if strings.Contains(err.Error(), "apply your changes to the latest version and try again") {
			log.Info("Conflict when updating Gateway listener status, retrying")
			return ctrl.Result{Requeue: true}, nil
		} else {
			log.Error(err, "Failed to update Gateway listener status")
			return ctrl.Result{}, err
		}
	}

	tunnels, err := api.ZeroTrust.Tunnels.List(ctx, zero_trust.TunnelListParams{
		AccountID: cloudflare.String(account),
		IsDeleted: cloudflare.Bool(false),
		Name:      cloudflare.String(gateway.Name),
	})
	if err != nil {
		if strings.Contains(err.Error(), "429 Too Many Requests") {
			log.Error(err, "Rate limited, requeueing after 5 minutes")
			return ctrl.Result{
				RequeueAfter: time.Minute * 5, // https://developers.cloudflare.com/fundamentals/api/reference/limits/
			}, nil
		} else {
			return ctrl.Result{}, err
		}
	}

	tunnelID := ""
	if len(tunnels.Result) == 0 {
		log.Info("Creating tunnel")
		// secret is required, despite optional in docs and seemingly only needed for ConfigSrc=local
		tunnel, err := api.ZeroTrust.Tunnels.New(ctx, zero_trust.TunnelNewParams{
			AccountID:    cloudflare.String(account),
			Name:         cloudflare.String(gateway.Name),
			TunnelSecret: cloudflare.String("AQIDBAUGBwgBAgMEBQYHCAECAwQFBgcIAQIDBAUGBwg="),
			// config_src = cloudflare
		})
		if err != nil {
			log.Error(err, "Failed to create tunnel")
			return ctrl.Result{}, err
		}
		tunnelID = tunnel.ID
	} else {
		// patch unsupported with api_token
		// if tunnels[0].Name != gateway.Name {
		// log.Info("updating Tunnel name")
		// API uses /cfd_tunnel/{id}, but SDK uses /cfd_tunnel? might be broken
		// _, err := api.UpdateTunnel(ctx, account, cloudflare.TunnelUpdateParams{Name: gateway.Name})
		// if err != nil {
		// 	log.Error(err, "unable to update Tunnel")
		// 	return ctrl.Result{}, err
		// }
		// }
		tunnelID = tunnels.Result[0].ID
	}

	res, err := api.ZeroTrust.Tunnels.Token.Get(ctx, tunnelID, zero_trust.TunnelTokenGetParams{
		AccountID: cloudflare.String(account),
	})
	if err != nil {
		log.Error(err, "Failed to get tunnel token")
		return ctrl.Result{}, err
	}
	token := string((*res).(shared.UnionString))

	// Check if the deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: gateway.Name, Namespace: gateway.Namespace}, found)
	// TODO update existing deployment eg image changes
	if err != nil && apierrors.IsNotFound(err) {
		// Define a new deployment
		dep, err := r.deploymentForGateway(gateway, token)
		if err != nil {
			log.Error(err, "Failed to define new Deployment resource for Gateway")

			// The following implementation will update the status
			meta.SetStatusCondition(&gateway.Status.Conditions, metav1.Condition{Type: string(gatewayv1.GatewayConditionAccepted),
				Status: metav1.ConditionFalse, Reason: "Reconciling", ObservedGeneration: gateway.Generation,
				Message: fmt.Sprintf("Failed to create Deployment for the custom resource (%s): (%s)", gateway.Name, err)})

			if err := r.Status().Update(ctx, gateway); err != nil {
				log.Error(err, "Failed to update Gateway status")
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, err
		}

		log.Info("Creating a new Deployment",
			"Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		if err = r.Create(ctx, dep); err != nil {
			log.Error(err, "Failed to create new Deployment",
				"Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}

		// Deployment created successfully
		// We will requeue the reconciliation so that we can ensure the state
		// and move forward for the next operations
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		// Let's return the error for the reconciliation be re-trigged again
		return ctrl.Result{}, err
	} else {
		// Define a new deployment
		dep, err := r.deploymentForGateway(gateway, token)
		if err != nil {
			log.Error(err, "Failed to define new Deployment resource for Gateway")

			// The following implementation will update the status
			meta.SetStatusCondition(&gateway.Status.Conditions, metav1.Condition{Type: string(gatewayv1.GatewayConditionAccepted),
				Status: metav1.ConditionFalse, Reason: "Reconciling", ObservedGeneration: gateway.Generation,
				Message: fmt.Sprintf("Failed to update Deployment for the custom resource (%s): (%s)", gateway.Name, err)})

			if err := r.Status().Update(ctx, gateway); err != nil {
				log.Error(err, "Failed to update Gateway status")
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, err
		}

		if err := r.Update(ctx, dep); err != nil {
			if strings.Contains(err.Error(), "apply your changes to the latest version and try again") {
				log.Info("Conflict when updating Deployment, retrying")
				return ctrl.Result{Requeue: true}, nil
			} else {
				log.Error(err, "Failed to update Deployment")
				return ctrl.Result{}, err
			}
		}
	}

	if err := r.Get(ctx, req.NamespacedName, gateway); err != nil {
		log.Error(err, "Failed to re-fetch gateway")
		return ctrl.Result{}, err
	}

	// The following implementation will update the status
	meta.SetStatusCondition(&gateway.Status.Conditions, metav1.Condition{Type: string(gatewayv1.GatewayConditionProgrammed),
		Status: metav1.ConditionTrue, Reason: string(gatewayv1.GatewayReasonProgrammed), ObservedGeneration: gateway.Generation,
		Message: fmt.Sprintf("Tunnel and deployment for gateway (%s) reconciled successfully", gateway.Name)})

	if err := r.Status().Update(ctx, gateway); err != nil {
		if strings.Contains(err.Error(), "apply your changes to the latest version and try again") {
			log.Info("Conflict when updating Gateway listener status, retrying")
			return ctrl.Result{Requeue: true}, nil
		} else {
			log.Error(err, "Failed to update Gateway status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// finalizeGateway will perform the required operations before delete the CR.
func (r *GatewayReconciler) doFinalizerOperationsForGateway(ctx context.Context, gatewayClass *gatewayv1.GatewayClass, gateway *gatewayv1.Gateway, account string, api *cloudflare.Client) error {
	// Note: It is not recommended to use finalizers with the purpose of deleting resources which are
	// created and managed in the reconciliation. These ones, such as the Deployment created on this reconcile,
	// are defined as dependent of the custom resource. See that we use the method ctrl.SetControllerReference.
	// to set the ownerRef which means that the Deployment will be deleted by the Kubernetes API.
	// More info: https://kubernetes.io/docs/tasks/administer-cluster/use-cascading-deletion/

	log := log.FromContext(ctx)

	tunnel, err := api.ZeroTrust.Tunnels.List(ctx, zero_trust.TunnelListParams{
		AccountID: cloudflare.String(account),
		IsDeleted: cloudflare.Bool(false),
		Name:      cloudflare.String(gateway.Name),
	})
	if err != nil {
		log.Error(err, "Failed to get tunnel from Cloudflare API")
		return err
	}

	if len(tunnel.Result) > 0 {
		log.Info("Deleting Tunnel")

		// prepare for deletion - disconnect, after rotating secret to prevent reconnect
		if _, err := api.ZeroTrust.Tunnels.Edit(ctx, tunnel.Result[0].ID, zero_trust.TunnelEditParams{
			AccountID:    cloudflare.String(account),
			TunnelSecret: cloudflare.String("Vm0xd1MwMUhSWGhYV0d4VlYwZG9jVlZ0TVRSV01XeHpZVWR3VUZWVU1Eaz0="),
		}); err != nil {
			log.Error(err, "Failed to update tunnel secret")
			return err
		}
		if _, err := api.ZeroTrust.Tunnels.Connections.Delete(ctx, tunnel.Result[0].ID, zero_trust.TunnelConnectionDeleteParams{
			AccountID: cloudflare.String(account),
		}); err != nil {
			log.Error(err, "Failed to delete tunnel connections")
			return err
		}

		if _, err := api.ZeroTrust.Tunnels.Delete(ctx, tunnel.Result[0].ID, zero_trust.TunnelDeleteParams{
			AccountID: cloudflare.String(account),
		}); err != nil {
			log.Error(err, "Failed to delete tunnel")
			return err
		}
	} else {
		log.Info("Gateway under deletion has no tunnel")
	}

	// if GatewayClass has no other Gateways, remove its finalizer
	gateways := &gatewayv1.GatewayList{Items: []gatewayv1.Gateway{{Spec: gatewayv1.GatewaySpec{GatewayClassName: gateway.Spec.GatewayClassName}}}}
	if err := r.List(ctx, gateways); err != nil {
		log.Error(err, "Failed to list Gateways")
		return err
	}
	if len(gateways.Items) == 0 {
		controllerutil.RemoveFinalizer(gatewayClass, gatewayClassFinalizer)
		if err := r.Update(ctx, gatewayClass); err != nil {
			return err
		}
	}

	// The following implementation will raise an event
	r.Recorder.Event(gateway, "Warning", "Deleting",
		fmt.Sprintf("Gateway %s is being deleted from the namespace %s",
			gateway.Name,
			gateway.Namespace))

	return nil
}

// deploymentForGateway returns a Gateway Deployment object
func (r *GatewayReconciler) deploymentForGateway(
	gateway *gatewayv1.Gateway, token string) (*appsv1.Deployment, error) {
	ls := labelsForGateway(gateway.Name)
	replicas := int32(1)

	// Get the Operand image
	image, err := imageForGateway()
	if err != nil {
		return nil, err
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gateway.Name,
			Namespace: gateway.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Affinity: &corev1.Affinity{
						NodeAffinity: &corev1.NodeAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
								NodeSelectorTerms: []corev1.NodeSelectorTerm{
									{
										MatchExpressions: []corev1.NodeSelectorRequirement{
											{
												Key:      "kubernetes.io/arch",
												Operator: "In",
												Values:   []string{"amd64", "arm64"},
											},
											{
												Key:      "kubernetes.io/os",
												Operator: "In",
												Values:   []string{"linux"},
											},
										},
									},
								},
							},
						},
					},
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot: &[]bool{true}[0],
						// IMPORTANT: seccomProfile was introduced with Kubernetes 1.19
						// If you are looking for to produce solutions to be supported
						// on lower versions you must remove this option.
						SeccompProfile: &corev1.SeccompProfile{
							Type: corev1.SeccompProfileTypeRuntimeDefault,
						},
					},
					Containers: []corev1.Container{{
						Image:           image,
						Name:            "gateway",
						ImagePullPolicy: corev1.PullIfNotPresent,
						// Ensure restrictive context for the container
						// More info: https://kubernetes.io/docs/concepts/security/pod-security-standards/#restricted
						SecurityContext: &corev1.SecurityContext{
							RunAsNonRoot:             &[]bool{true}[0],
							RunAsUser:                &[]int64{1001}[0],
							AllowPrivilegeEscalation: &[]bool{false}[0],
							Capabilities: &corev1.Capabilities{
								Drop: []corev1.Capability{
									"ALL",
								},
							},
						},
						Args: []string{"tunnel", "--no-autoupdate", "--metrics", "0.0.0.0:2000", "run", "--token", token},
					}},
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{IntVal: 0},
				},
			},
		},
	}

	// Set the ownerRef for the Deployment
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/owners-dependents/
	if err := ctrl.SetControllerReference(gateway, dep, r.Scheme); err != nil {
		return nil, err
	}
	return dep, nil
}

// labelsForGateway returns the labels for selecting the resources
// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
func labelsForGateway(name string) map[string]string {
	// skip imageTag, to allow version updates to existing deployments
	// var imageTag string
	// image, err := imageForGateway()
	// if err == nil {
	// 	imageTag = strings.Split(image, ":")[1]
	// }
	return map[string]string{
		"app.kubernetes.io/name": "cloudflare-kubernetes-gateway",
		// "app.kubernetes.io/version":    imageTag,
		"app.kubernetes.io/managed-by": "GatewayController",
		"cfargotunnel.com/name":        name,
	}
}

// imageForGateway gets the Operand image which is managed by this controller
// from the GATEWAY_IMAGE environment variable defined in the config/manager/manager.yaml
func imageForGateway() (string, error) {
	var imageEnvVar = "GATEWAY_IMAGE"
	image, found := os.LookupEnv(imageEnvVar)
	if !found {
		return "", fmt.Errorf("unable to find %s environment variable with the image", imageEnvVar)
	}
	return image, nil
}

// SetupWithManager sets up the controller with the Manager.
// Note that the Deployment will be also watched in order to ensure its
// desirable state on the cluster
func (r *GatewayReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gatewayv1.Gateway{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
