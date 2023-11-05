package controller

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	gw "sigs.k8s.io/gateway-api/apis/v1"
)

// HTTPRouteReconciler reconciles a HTTPRoute object
type HTTPRouteReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	Namespace string
}

//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=httproutes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=httproutes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=httproutes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// This is very inefficient, especially the API calls. But it shouldn't matter for small deployments.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile

// TODO handle deletion? eg handle all HTTPRoutes at once. how to reconcile DNS though? tunnel ID in comment field?
func (r *HTTPRouteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	route := &gw.HTTPRoute{}
	if err := r.Get(ctx, req.NamespacedName, route); err != nil {
		// log.Error(err, "Failed to get HTTPRoute")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	sourceIngress := map[string]bool{}
	for _, rule := range route.Spec.Rules {
		paths := map[string]bool{}
		for _, match := range rule.Matches {
			if match.Path == nil {
				paths["/"] = true
			} else {
				paths[*match.Path.Value] = true
			}

			if match.Headers != nil {
				log.Info("HTTPRoute header match is not supported", match.Headers)
			}
		}

		// TODO implement this with rewrite rules? Core filters are a MUST in the spec
		if rule.Filters != nil {
			log.Info("HTTPRoute filters are not supported", rule.Filters)
		}

		services := map[string]bool{}
		for _, backend := range rule.BackendRefs {
			if backend.Port == nil {
				err := errors.New("HTTPRoute backend port is nil")
				log.Error(err, "HTTPRoute backend port is required and nil", backend)
				continue
			}

			var namespace string
			if backend.Namespace == nil {
				namespace = route.Namespace
			} else {
				namespace = string(*backend.Namespace)
			}

			services[fmt.Sprintf("http://%s.%s:%d", string(backend.Name), namespace, int32(*backend.Port))] = true
		}

		for _, hostname := range route.Spec.Hostnames {
			for path := range paths {
				for service := range services {
					sourceIngress[fmt.Sprintf("%s|%s|%s", hostname, path, service)] = true
				}
			}
		}
	}

	for _, parentRef := range route.Spec.ParentRefs {
		gateway := &gw.Gateway{}
		if err := r.Get(ctx, types.NamespacedName{
			Namespace: string(*parentRef.Namespace),
			Name:      string(parentRef.Name),
		}, gateway); err != nil {
			log.Error(err, "Failed to get Gateway")
			return ctrl.Result{}, err
		}

		account, api, err := InitCloudflareApi(ctx, r.Client, string(gateway.Spec.GatewayClassName))
		if err != nil {
			log.Error(err, "Failed to initialize Cloudflare API")
			return ctrl.Result{}, err
		}

		tunnels, _, err := api.ListTunnels(ctx, account, cloudflare.TunnelListParams{IsDeleted: cloudflare.BoolPtr(false), Name: gateway.Name})
		if err != nil {
			log.Error(err, "Failed to get tunnel from Cloudflare API")
			return ctrl.Result{}, err
		}
		tunnel := tunnels[0]

		config, err := api.GetTunnelConfiguration(ctx, account, tunnel.ID)
		if err != nil {
			log.Error(err, "Failed to get tunnel config from Cloudflare API")
			return ctrl.Result{}, err
		}

		ingress := sourceIngress
		for _, rule := range config.Config.Ingress {
			// skip catch-all rules
			if rule.Hostname != "" {
				ingress[fmt.Sprintf("%s|%s|%s", rule.Hostname, rule.Path, rule.Service)] = true
			}
		}

		ingressStruct := []cloudflare.UnvalidatedIngressRule{}
		for rule := range ingress {
			parts := strings.Split(rule, "|")
			ingressStruct = append(ingressStruct, cloudflare.UnvalidatedIngressRule{
				Hostname: parts[0],
				Path:     parts[1],
				Service:  parts[2],
			})
		}

		// last rule must be the catch-all
		ingressStruct = append(ingressStruct, cloudflare.UnvalidatedIngressRule{
			Service: "http_status:404",
		})

		_, err = api.UpdateTunnelConfiguration(ctx, account, cloudflare.TunnelConfigurationParams{
			TunnelID: tunnels[0].ID,
			Config: cloudflare.TunnelConfiguration{
				Ingress: ingressStruct,
			},
		})
		if err != nil {
			log.Error(err, "Failed to update Tunnel configuration")
			return ctrl.Result{}, err
		}

		log.Info("Updated Tunnel configuration", "ingress", ingressStruct)

		for _, gwHostname := range route.Spec.Hostnames {
			hostname := string(gwHostname)
			// terrible
			zoneName := strings.Join(strings.Split(hostname, ".")[1:], ".")

			zones, err := api.ListZonesContext(ctx, cloudflare.WithZoneFilters(zoneName, account.Identifier, "active"))
			if err != nil {
				log.Error(err, "Failed to list DNS zones")
				return ctrl.Result{}, err
			}
			zone := cloudflare.ResourceIdentifier(zones.Result[0].ID)

			content := fmt.Sprintf("%s.cfargotunnel.com", tunnel.ID)
			records, info, _ := api.ListDNSRecords(ctx, zone, cloudflare.ListDNSRecordsParams{
				Proxied: cloudflare.BoolPtr(true),
				Type:    "CNAME",
				Name:    hostname,
				Content: content,
			})
			if info.Count == 0 {
				_, err := api.CreateDNSRecord(ctx, zone, cloudflare.CreateDNSRecordParams{
					Proxied: cloudflare.BoolPtr(true),
					Type:    "CNAME",
					Name:    hostname,
					Content: content,
				})
				if err != nil {
					log.Error(err, "Failed to create DNS record", hostname, content)
					return ctrl.Result{}, err
				}
			} else {
				_, err := api.UpdateDNSRecord(ctx, zone, cloudflare.UpdateDNSRecordParams{
					ID:      records[0].ID,
					Proxied: cloudflare.BoolPtr(true),
					Type:    "CNAME",
					Name:    hostname,
					Content: content,
				})
				if err != nil {
					log.Error(err, "Failed to update DNS record", hostname, content)
					return ctrl.Result{}, err
				}
			}
		}
		log.Info("Updated DNS records", "hostnames", route.Spec.Hostnames)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HTTPRouteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gw.HTTPRoute{}).
		Complete(r)
}
