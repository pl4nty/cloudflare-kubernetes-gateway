package controller

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/accounts"
	"github.com/cloudflare/cloudflare-go/v7/option"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	gw "sigs.k8s.io/gateway-api/apis/v1"
)

func InitCloudflareAPI(ctx context.Context, c client.Client, gatewayClassName string) (string, *cloudflare.Client, error) {
	log := log.FromContext(ctx)

	gatewayClass := &gw.GatewayClass{}
	if err := c.Get(ctx, types.NamespacedName{Name: gatewayClassName}, gatewayClass); err != nil {
		log.Error(err, "Failed to get gatewayclass")
		return "", nil, err
	}
	if gatewayClass.Spec.ControllerName != "github.com/pl4nty/cloudflare-kubernetes-gateway" {
		return "", nil, nil
	}

	if gatewayClass.Spec.ParametersRef == nil {
		return "", nil, errors.New("GatewayClass is missing a Secret ParameterRef")
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

func VerifyAPIToken(ctx context.Context, account string, api *cloudflare.Client) (string, error) {
	token, err := api.User.Tokens.Verify(ctx)
	if err != nil {
		token, err := api.Accounts.Tokens.Verify(ctx, accounts.TokenVerifyParams{AccountID: cloudflare.String(account)})
		if err != nil {
			return err.Error() + " Ensure ACCOUNT_ID and TOKEN are valid", nil
		} else {
			if token.Status != "active" {
				return fmt.Sprintf("Token status is %s, not active. Please check the Cloudflare dashboard", token.Status), nil
			}
		}
	} else {
		if token.Status != "active" {
			return fmt.Sprintf("Token status is %s, not active. Please check the Cloudflare dashboard", token.Status), nil
		}
	}
	return "", nil
}
