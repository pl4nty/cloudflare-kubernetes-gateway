package controller

import (
	"context"
	"errors"

	"github.com/cloudflare/cloudflare-go"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	gw "sigs.k8s.io/gateway-api/apis/v1"
)

// GatewayReconciler reconciles a Gateway object
type GatewayReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	Namespace string
}

//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gatewayclasses,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO implement actual Gateway API spec, eg Status field
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *GatewayReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	gateway := &gw.Gateway{}
	if err := r.Get(ctx, req.NamespacedName, gateway); err != nil {
		log.Error(err, "Failed to get Gateway")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// check if parent GatewayClass is ours and update finalizer
	gatewayClass := &gw.GatewayClass{}
	if err := r.Get(ctx, types.NamespacedName{Name: string(gateway.Spec.GatewayClassName)}, gatewayClass); err != nil {
		log.Error(err, "Failed to get GatewayClass")
		return ctrl.Result{}, err
	}
	if gatewayClass.Spec.ControllerName != "github.com/pl4nty/cloudflare-kubernetes-gateway" {
		return ctrl.Result{}, nil
	}
	gatewayClassFinalizer := "gateway-exists-finalizer.gateway.networking.k8s.io"
	if !controllerutil.ContainsFinalizer(gatewayClass, gatewayClassFinalizer) {
		controllerutil.AddFinalizer(gatewayClass, gatewayClassFinalizer)
		if err := r.Update(ctx, gatewayClass); err != nil {
			return ctrl.Result{}, err
		}
	}

	// check spec requirement for at least one listener
	if len(gateway.Spec.Listeners) == 0 {
		err := errors.New("invalid spec.listeners")
		log.Error(err, "Invalid Gateway spec.listeners, at least one listener must be specified")
		return ctrl.Result{}, err
	}

	// TODO Gateway status
	// names := &map[string]{}
	// for _, listener := range gateway.Spec.Listeners {
	// 	value, ok := names[string(listener.Name)]

	// 	gatewayClass.Status.Conditions = append(gatewayClass.Status.Conditions, condition)
	// 	if err := r.Update(ctx, gatewayClass); err != nil {
	// 		log.Error(err, "Failed to update GatewayClass status")
	// 		return ctrl.Result{}, err
	// 	}
	// }

	account, api, err := InitCloudflareApi(ctx, r.Client, gatewayClass.Name)
	if err != nil {
		log.Error(err, "Failed to initialize Cloudflare API")
		return ctrl.Result{}, err
	}

	// reconcile Gateway finalizer
	gatewayFinalizer := "cfargotunnel.com/finalizer"
	if gateway.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(gateway, gatewayFinalizer) {
			controllerutil.AddFinalizer(gateway, gatewayFinalizer)
			if err := r.Update(ctx, gateway); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(gateway, gatewayFinalizer) {
			// TODO better identifiers
			if err := r.Delete(ctx, &apps.Deployment{
				ObjectMeta: v1.ObjectMeta{
					Namespace: gateway.Namespace,
					Name:      gateway.Name,
				},
			}); err != nil {
				log.Error(err, "Failed to delete Deployment")
			}

			tunnel, tunnelInfo, err := api.ListTunnels(ctx, account, cloudflare.TunnelListParams{
				IsDeleted: cloudflare.BoolPtr(false),
				Name:      gateway.Name,
			})
			if err != nil {
				log.Error(err, "Failed to get tunnel from Cloudflare API")
				return ctrl.Result{}, err
			}

			if tunnelInfo.Count > 0 {
				log.Info("Deleting Tunnel")
				if err := api.DeleteTunnel(ctx, account, tunnel[0].ID); err != nil {
					log.Error(err, "Failed to delete tunnel Deployment")
					return ctrl.Result{}, err
				}
			} else {
				log.Info("Gateway under deletion has no tunnel")
			}

			controllerutil.RemoveFinalizer(gateway, gatewayFinalizer)
			if err := r.Update(ctx, gateway); err != nil {
				return ctrl.Result{}, err
			}
		}

		// if GatewayClass has no other Gateways, remove its finalizer
		gateways := &gw.GatewayList{Items: []gw.Gateway{{Spec: gw.GatewaySpec{GatewayClassName: gateway.Spec.GatewayClassName}}}}
		if err := r.List(ctx, gateways); err != nil {
			log.Error(err, "Failed to list Gateways")
			return ctrl.Result{}, err
		}
		if len(gateways.Items) == 0 {
			controllerutil.RemoveFinalizer(gatewayClass, gatewayClassFinalizer)
			if err := r.Update(ctx, gatewayClass); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}

	tunnels, info, err := api.ListTunnels(ctx, account, cloudflare.TunnelListParams{IsDeleted: cloudflare.BoolPtr(false), Name: gateway.Name})
	if err != nil {
		log.Error(err, "Failed to get Tunnel from Cloudflare API")
		return ctrl.Result{}, err
	}

	tunnel := cloudflare.Tunnel{}
	if info.Count == 0 {
		log.Info("Creating tunnel")
		// secret is required, despite optional in docs and seemingly only needed for ConfigSrc=local
		tunnel, err = api.CreateTunnel(ctx, account, cloudflare.TunnelCreateParams{
			Name:      gateway.Name,
			ConfigSrc: "cloudflare",
			Secret:    "AQIDBAUGBwgBAgMEBQYHCAECAwQFBgcIAQIDBAUGBwg=",
		})
		if err != nil {
			log.Error(err, "Failed to create tunnel")
			return ctrl.Result{}, err
		}
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
		log.Info("Tunnel exists")
		tunnel = tunnels[0]
	}

	token, err := api.GetTunnelToken(ctx, account, tunnel.ID)
	if err != nil {
		log.Error(err, "Failed to get tunnel token")
		return ctrl.Result{}, err
	}

	if err := r.Get(ctx, types.NamespacedName{
		Namespace: gateway.Namespace,
		Name:      gateway.Name,
	}, &apps.Deployment{}); err == nil {
		log.Info("Tunnel deployment exists")
		return ctrl.Result{}, nil
	}

	labels := map[string]string{"cfargotunnel.com/name": gateway.Name}
	deployment := apps.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Namespace: gateway.Namespace,
			Name:      gateway.Name,
		},
		Spec: apps.DeploymentSpec{
			Selector: &v1.LabelSelector{MatchLabels: labels},
			Template: core.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{Labels: labels},
				Spec: core.PodSpec{Containers: []core.Container{{
					Name:  "main",
					Image: "cloudflare/cloudflared:2023.8.2",
					Args:  []string{"tunnel", "--no-autoupdate", "run", "--token", token},
				}}},
			},
		},
	}

	if err := r.Create(ctx, &deployment); err != nil {
		log.Error(err, "Failed to create tunnel deployment")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GatewayReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gw.Gateway{}).
		Owns(&apps.Deployment{}).
		Complete(r)
}
