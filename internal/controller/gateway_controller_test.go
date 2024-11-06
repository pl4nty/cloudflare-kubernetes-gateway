package controller

import (
	"context"
	"os"
	"time"

	//nolint:golint
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

var _ = Describe("Gateway controller", func() {
	Context("Gateway controller test", func() {

		const GatewayName = "test-gateway"

		ctx := context.Background()

		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      GatewayName,
				Namespace: GatewayName,
			},
		}

		typeNamespaceName := types.NamespacedName{
			Name:      GatewayName,
			Namespace: GatewayName,
		}
		gateway := &gatewayv1.Gateway{}

		BeforeEach(func() {
			By("Creating the Namespace to perform the tests")
			err := k8sClient.Create(ctx, namespace)
			Expect(err).To(Not(HaveOccurred()))

			By("Setting the Image ENV VAR which stores the Operand image")
			// renovate: datasource=docker depName=cloudflare/cloudflared
			err = os.Setenv("GATEWAY_IMAGE", "cloudflare/cloudflared:2024.11.0")
			Expect(err).To(Not(HaveOccurred()))

			By("creating the custom resource for the Kind Gateway")
			err = k8sClient.Get(ctx, typeNamespaceName, gateway)
			if err != nil && errors.IsNotFound(err) {
				// Let's mock our custom resource at the same way that we would
				// apply on the cluster the manifest under config/samples
				gateway := &gatewayv1.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Name:      GatewayName,
						Namespace: namespace.Name,
					},
					Spec: gatewayv1.GatewaySpec{
						GatewayClassName: "cloudflare",
						Listeners: []gatewayv1.Listener{{
							Name:     "cloudflare",
							Port:     80,
							Protocol: "http",
						}},
					},
				}

				err = k8sClient.Create(ctx, gateway)
				Expect(err).To(Not(HaveOccurred()))
			}
		})

		AfterEach(func() {
			By("removing the custom resource for the Kind Gateway")
			found := &gatewayv1.Gateway{}
			err := k8sClient.Get(ctx, typeNamespaceName, found)
			Expect(err).To(Not(HaveOccurred()))

			Eventually(func() error {
				return k8sClient.Delete(context.TODO(), found)
			}, 2*time.Minute, time.Second).Should(Succeed())

			// TODO(user): Attention if you improve this code by adding other context test you MUST
			// be aware of the current delete namespace limitations.
			// More info: https://book.kubebuilder.io/reference/envtest.html#testing-considerations
			By("Deleting the Namespace to perform the tests")
			_ = k8sClient.Delete(ctx, namespace)

			By("Removing the Image ENV VAR which stores the Operand image")
			_ = os.Unsetenv("GATEWAY_IMAGE")
		})

		It("should successfully reconcile a custom resource for Gateway", func() {
			By("Checking if the custom resource was successfully created")
			Eventually(func() error {
				found := &gatewayv1.Gateway{}
				return k8sClient.Get(ctx, typeNamespaceName, found)
			}, time.Minute, time.Second).Should(Succeed())

			By("Reconciling the custom resource created")
			gatewayReconciler := &GatewayReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := gatewayReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespaceName,
			})
			Expect(err).To(Not(HaveOccurred()))

			// By("Checking if Deployment was successfully created in the reconciliation")
			// Eventually(func() error {
			// 	found := &appsv1.Deployment{}
			// 	return k8sClient.Get(ctx, typeNamespaceName, found)
			// }, time.Minute, time.Second).Should(Succeed())

			// By("Checking the latest Status Condition added to the Gateway instance")
			// Eventually(func() error {
			// 	if gateway.Status.Conditions != nil &&
			// 		len(gateway.Status.Conditions) != 0 {
			// 		latestStatusCondition := gateway.Status.Conditions[len(gateway.Status.Conditions)-1]
			// 		expectedLatestStatusCondition := metav1.Condition{
			// 			Type:   string(gatewayv1.GatewayConditionAccepted),
			// 			Status: metav1.ConditionTrue,
			// 			Reason: "Reconciling",
			// 			Message: fmt.Sprintf(
			// 				"Deployment for Gateway (%s) created successfully",
			// 				gateway.Name),
			// 		}
			// 		if latestStatusCondition != expectedLatestStatusCondition {
			// 			return fmt.Errorf("The latest status condition added to the Gateway instance is not as expected")
			// 		}
			// 	}
			// 	return nil
			// }, time.Minute, time.Second).Should(Succeed())
		})
	})
})
