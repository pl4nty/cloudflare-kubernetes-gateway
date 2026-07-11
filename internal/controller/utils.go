package controller

import (
	"context"
	"errors"
	"strings"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/option"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	gw "sigs.k8s.io/gateway-api/apis/v1"
)

// const controllerName declared in gateway_controller.go

func InitCloudflareAPI(ctx context.Context, c client.Client, gatewayClassName string) (string, *cloudflare.Client, error) {
	accountID, apiToken, err := GetCloudflareAPICredentials(ctx, c, gatewayClassName)
	if err != nil {
		return "", nil, err
	}
	api := cloudflare.NewClient(option.WithAPIToken(apiToken))
	return accountID, api, nil
}

func GetCloudflareAPICredentials(ctx context.Context, c client.Client, gatewayClassName string) (string, string, error) {
	logger := log.FromContext(ctx)

	gatewayClass := &gw.GatewayClass{}
	if err := c.Get(ctx, types.NamespacedName{Name: gatewayClassName}, gatewayClass); err != nil {
		logger.Error(err, "Failed to get gatewayclass")
		return "", "", err
	}
	if gatewayClass.Spec.ControllerName != controllerName {
		return "", "", nil
	}

	if gatewayClass.Spec.ParametersRef == nil {
		return "", "", errors.New("GatewayClass is missing a Secret ParameterRef")
	}

	secretRef := types.NamespacedName{
		Namespace: string(*gatewayClass.Spec.ParametersRef.Namespace),
		Name:      gatewayClass.Spec.ParametersRef.Name,
	}
	secret := &core.Secret{}
	if err := c.Get(ctx, secretRef, secret); err != nil {
		logger.Error(err, "unable to fetch Secret from GatewayClass ParameterRef")
		return "", "", err
	}

	accountID := strings.TrimSpace(string(secret.Data["ACCOUNT_ID"]))
	apiToken := strings.TrimSpace(string(secret.Data["TOKEN"]))
	return accountID, apiToken, nil
}
