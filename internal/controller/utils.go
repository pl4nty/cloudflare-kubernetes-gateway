package controller

import (
	"context"
	"strings"

	"github.com/cloudflare/cloudflare-go/v2"
	"github.com/cloudflare/cloudflare-go/v2/option"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	gw "sigs.k8s.io/gateway-api/apis/v1"
)

func InitCloudflareApi(ctx context.Context, c client.Client, gatewayClassName string) (string, *cloudflare.Client, error) {
	log := log.FromContext(ctx)

	gatewayClass := &gw.GatewayClass{}
	if err := c.Get(ctx, types.NamespacedName{Name: gatewayClassName}, gatewayClass); err != nil {
		log.Error(err, "Failed to get gatewayclass")
		return "", nil, err
	}
	if gatewayClass.Spec.ControllerName != "github.com/pl4nty/cloudflare-kubernetes-gateway" {
		return "", nil, nil
	}

	secretRef := types.NamespacedName{
		Namespace: string(*gatewayClass.Spec.ParametersRef.Namespace),
		Name:      gatewayClass.Spec.ParametersRef.Name,
	}
	secret := &core.Secret{}
	if err := c.Get(ctx, secretRef, secret); err != nil {
		log.Error(err, "unable to fetch Secret from GatewayClass ParameterRef")
		return "", nil, err
	}

	account := strings.TrimSpace(string(secret.Data["ACCOUNT_ID"]))
	api := cloudflare.NewClient(option.WithAPIToken(strings.TrimSpace(string(secret.Data["TOKEN"]))))

	return account, api, nil
}
