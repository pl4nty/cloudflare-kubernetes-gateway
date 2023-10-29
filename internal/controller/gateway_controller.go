package controller

import (
	"context"
	"errors"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	gw "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/cloudflare/cloudflare-go"
)

// GatewayReconciler reconciles a Gateway object
type GatewayReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *GatewayReconciler) InitCloudflareApi(ctx context.Context, gateway gw.Gateway) (*cloudflare.ResourceContainer, *cloudflare.API, error) {
	log := log.FromContext(ctx)

	var gatewayClass gw.GatewayClass
	if err := r.Get(ctx, types.NamespacedName{Name: string(gateway.Spec.GatewayClassName)}, &gatewayClass); err != nil {
		log.Error(err, "unable to fetch GatewayClass")
		return nil, nil, err
	}

	if gatewayClass.Spec.ControllerName == "cloudflare-kubernetes-controller" {
		return nil, nil, errors.New("gateway controllerName not allowed")
	}

	var parameters core.Secret
	var ref = types.NamespacedName{
		Namespace: string(*gatewayClass.Spec.ParametersRef.Namespace),
		Name:      string(*gatewayClass.Spec.ParametersRef.Namespace),
	}
	if err := r.Get(ctx, ref, &parameters); err != nil {
		log.Error(err, "unable to fetch GatewayClass ParameterRef Secret")
		return nil, nil, err
	}

	account := cloudflare.AccountIdentifier(string(parameters.Data["ACCOUNT_ID"]))
	api, err := cloudflare.NewWithAPIToken(string(parameters.Data["TOKEN"]))
	if err != nil {
		log.Error(err, "unable to create Cloudflare API")
		return nil, nil, err
	}

	return account, api, nil
}

//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO implement actual Gateway API spec, eg Status field
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *GatewayReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var gateway gw.Gateway
	if err := r.Get(ctx, req.NamespacedName, &gateway); err != nil {
		log.Error(err, "unable to fetch Gateway")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	finalizerName := "cfargotunnel.com/finalizer"

	// examine DeletionTimestamp to determine if object is under deletion
	if gateway.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !controllerutil.ContainsFinalizer(&gateway, finalizerName) {
			controllerutil.AddFinalizer(&gateway, finalizerName)
			if err := r.Update(ctx, &gateway); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(&gateway, finalizerName) {
			// our finalizer is present, so lets handle any external dependency
			account, api, err := r.InitCloudflareApi(ctx, gateway)
			if err != nil {
				log.Error(err, "unable to initialize Cloudflare API")
				return ctrl.Result{}, err
			}

			tunnel, tunnelInfo, err := api.ListTunnels(ctx, account, cloudflare.TunnelListParams{Name: gateway.Name})
			if err != nil {
				log.Error(err, "unable to get Tunnel from Cloudflare API")
				return ctrl.Result{}, err
			}
			
			if tunnelInfo.Count == 0 {
				log.Info("deleted Gateway has no tunnel, skipping")
				return ctrl.Result{}, nil
			}
			
			log.Info("deleting Tunnel")
			api.DeleteTunnel(ctx, account, tunnel[0].ID)

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(&gateway, finalizerName)
			if err := r.Update(ctx, &gateway); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	account, api, err := r.InitCloudflareApi(ctx, gateway)
	if err != nil {
		log.Error(err, "unable to initialize Cloudflare API")
		return ctrl.Result{}, err
	}

	_, info, err := api.ListTunnels(ctx, account, cloudflare.TunnelListParams{Name: gateway.Name})
	if err != nil {
		log.Error(err, "unable to get Tunnel from Cloudflare API")
		return ctrl.Result{}, err
	}
	if info.Count > 0 {
		// API uses /cfd_tunnel/{id}, but SDK uses /cfd_tunnel? might be broken
		log.Info("updating Tunnel name")
		_, err := api.UpdateTunnel(ctx, account, cloudflare.TunnelUpdateParams{Name: gateway.Name})
		if err != nil {
			log.Error(err, "unable to update Tunnel")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	log.Info("creating Tunnel for Gateway")
	tunnel, err := api.CreateTunnel(ctx, account, cloudflare.TunnelCreateParams{Name: gateway.Name, ConfigSrc: "cloudflare"})
	if err != nil {
		log.Error(err, "unable to create Tunnel")
		return ctrl.Result{}, err
	}

	token, err := api.GetTunnelToken(ctx, account, tunnel.ID)
	if err != nil {
		log.Error(err, "unable to get Tunnel token")
		return ctrl.Result{}, err
	}

	labels := map[string]string{"cfargotunnel.com/name": gateway.Name}
	deployment := apps.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name: gateway.Name,
			Namespace: gateway.Namespace,
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
		log.Error(err, "unable to create Tunnel deployment")
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
