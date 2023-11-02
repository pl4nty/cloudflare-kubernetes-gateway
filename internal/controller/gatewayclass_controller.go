package controller

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	gw "sigs.k8s.io/gateway-api/apis/v1"
)

// GatewayClassReconciler reconciles a GatewayClass object
type GatewayClassReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gatewayclasses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gatewayclasses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gatewayclasses/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GatewayClass object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *GatewayClassReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	gatewayClass := &gw.GatewayClass{}
	if err := r.Get(ctx, req.NamespacedName, gatewayClass); err != nil {
		log.Error(err, "unable to fetch GatewayClass")
		return ctrl.Result{}, err
	}

	// TODO is this set by default?
	gatewayClass.Status.Conditions = append(gatewayClass.Status.Conditions, v1.Condition{
		Type:               "Accepted",
		Status:             v1.ConditionUnknown,
		Reason:             "Pending",
		LastTransitionTime: v1.Now(),
		ObservedGeneration: gatewayClass.Generation,
	})
	if err := r.Update(ctx, gatewayClass); err != nil {
		log.Error(err, "unable to update GatewayClass status")
		return ctrl.Result{}, err
	}

	// validate parameters
	condition := v1.Condition{}
	_, _, err := InitCloudflareApi(ctx, r.Client, gatewayClass.Name)
	if err != nil {
		log.Error(err, "unable to initialize Cloudflare API from secret in GatewayClass parameterRef. Ensure ACCOUNT_ID and TOKEN are set")

		condition = v1.Condition{
			Type:               "Accepted",
			Status:             v1.ConditionFalse,
			Reason:             "InvalidParameters",
			Message:            "Unable to initialize Cloudflare API from secret in GatewayClass parameterRef. Ensure ACCOUNT_ID and TOKEN are set",
			LastTransitionTime: v1.Now(),
			ObservedGeneration: gatewayClass.Generation,
		}
	} else {
		condition = v1.Condition{
			Type:               "Accepted",
			Status:             v1.ConditionTrue,
			Reason:             "Accepted",
			LastTransitionTime: v1.Now(),
			ObservedGeneration: gatewayClass.Generation,
		}
	}

	gatewayClass.Status.Conditions = append(gatewayClass.Status.Conditions, condition)
	if err := r.Update(ctx, gatewayClass); err != nil {
		log.Error(err, "unable to update GatewayClass status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GatewayClassReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gw.GatewayClass{Spec: gw.GatewayClassSpec{ControllerName: "github.com/pl4nty/cloudflare-kubernetes-controller"}}).
		Complete(r)
}
