package controller

import (
	"context"
	"fmt"

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
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HTTPRoute object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *HTTPRouteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	route := &gw.HTTPRoute{}
	if err := r.Get(ctx, req.NamespacedName, route); err != nil {
		log.Error(err, "unable to fetch HTTPRoute")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	for _, parentRef := range route.Spec.ParentRefs {
		gateway := &gw.Gateway{}
		if err := r.Get(ctx, types.NamespacedName{
			Namespace: string(*parentRef.Namespace),
			Name:      string(parentRef.Name),
		}, gateway); err != nil {
			log.Error(err, "unable to fetch Gateway")
			return ctrl.Result{}, err
		}

		account, api, err := InitCloudflareApi(ctx, r.Client, string(gateway.Spec.GatewayClassName))
		if err != nil {
			log.Error(err, "unable to initialize Cloudflare API")
			return ctrl.Result{}, err
		}

		tunnels, _, err := api.ListTunnels(ctx, account, cloudflare.TunnelListParams{IsDeleted: cloudflare.BoolPtr(false), Name: gateway.Name})
		if err != nil {
			log.Error(err, "unable to get Tunnel from Cloudflare API")
			return ctrl.Result{}, err
		}
		// tunnel := tunnels[0]

		backendRef := route.Spec.Rules[0].BackendRefs[0]
		ingress := []cloudflare.UnvalidatedIngressRule{
			{
				Hostname: string(route.Spec.Hostnames[0]),
				Service:  fmt.Sprintf("http://%s.%s:%d", string(backendRef.Name), string(*backendRef.Namespace), int32(*backendRef.Port)),
			},
		}

		routes := &gw.HTTPRouteList{}
		if err := r.List(ctx, routes); err != nil {
			log.Error(err, "unable to fetch HTTPRoutes")
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}

		for _, anyRoute := range routes.Items {
			for _, anyParent := range anyRoute.Spec.ParentRefs {
				if anyParent == parentRef {
					hostname := string(anyRoute.Spec.Hostnames[0])
					anyBackendRef := anyRoute.Spec.Rules[0].BackendRefs[0]
					ingress = append(ingress, cloudflare.UnvalidatedIngressRule{
						Hostname: hostname,
						Service:  fmt.Sprintf("http://%s.%s:%d", string(anyBackendRef.Name), string(*anyBackendRef.Namespace), int32(*anyBackendRef.Port)),
					})

					// these need zone ID instead of account ID
					// lookup by hostname substring?
					// content := fmt.Sprintf("%s.cfargotunnel.com", tunnel.ID)
					// _, info, _ := api.ListDNSRecords(ctx, account, cloudflare.ListDNSRecordsParams{
					// 	Proxied: cloudflare.BoolPtr(true),
					// 	Type: "CNAME",
					// 	Name: hostname,
					// 	Content: content,
					// })
					// if info.Count == 0 {
					// 	api.CreateDNSRecord(ctx, account, cloudflare.CreateDNSRecordParams{
					// 		Proxied: cloudflare.BoolPtr(true),
					// 		Type: "CNAME",
					// 		Name: hostname,
					// 		Content: content,
					// 	})
					// } else {
					// 	api.UpdateDNSRecord(ctx, account, cloudflare.UpdateDNSRecordParams{
					// 		Proxied: cloudflare.BoolPtr(true),
					// 		Type: "CNAME",
					// 		Name: hostname,
					// 		Content: content,
					// 	})
					// }
				}
			}
		}

		// last rule must be the catch-all
		ingress = append(ingress, cloudflare.UnvalidatedIngressRule{
			Service: "http_status:404",
		})

		_, err = api.UpdateTunnelConfiguration(ctx, account, cloudflare.TunnelConfigurationParams{
			TunnelID: tunnels[0].ID,
			Config: cloudflare.TunnelConfiguration{
				Ingress: ingress,
			},
		})
		if err != nil {
			log.Error(err, "unable to update Tunnel configuration via Cloudflare API")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HTTPRouteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gw.HTTPRoute{}).
		Complete(r)
}
