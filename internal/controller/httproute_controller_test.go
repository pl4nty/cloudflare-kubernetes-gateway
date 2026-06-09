package controller

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

var _ = Describe("HTTPRoute Controller", func() {
	Context("When reconciling a resource", func() {

		It("should successfully reconcile the resource", func() {

			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})

	Context("When updating parent status", func() {
		It("sets Accepted for this controller and preserves other controller statuses", func() {
			namespace := gatewayv1.Namespace("default")
			parentRef := gatewayv1.ParentReference{
				Namespace: &namespace,
				Name:      gatewayv1.ObjectName("gateway"),
			}
			otherParentRef := gatewayv1.ParentReference{
				Namespace: &namespace,
				Name:      gatewayv1.ObjectName("other-gateway"),
			}
			route := &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Generation: 3,
				},
				Status: gatewayv1.HTTPRouteStatus{
					RouteStatus: gatewayv1.RouteStatus{
						Parents: []gatewayv1.RouteParentStatus{{
							ParentRef:      otherParentRef,
							ControllerName: gatewayv1.GatewayController("example.com/other-controller"),
							Conditions: []metav1.Condition{{
								Type:               string(gatewayv1.RouteConditionAccepted),
								Status:             metav1.ConditionFalse,
								Reason:             string(gatewayv1.RouteReasonNoMatchingParent),
								ObservedGeneration: 2,
								Message:            "owned by a different controller",
							}},
						}},
					},
				},
			}

			setHTTPRouteParentStatusCondition(route, parentRef, metav1.Condition{
				Type:               string(gatewayv1.RouteConditionAccepted),
				Status:             metav1.ConditionTrue,
				Reason:             string(gatewayv1.RouteReasonAccepted),
				ObservedGeneration: route.Generation,
				Message:            "Successfully reconciled with Cloudflare",
			})

			Expect(route.Status.Parents).To(HaveLen(2))
			Expect(route.Status.Parents[0].ControllerName).To(Equal(gatewayv1.GatewayController("example.com/other-controller")))

			parentStatus := route.Status.Parents[1]
			Expect(parentStatus.ParentRef).To(Equal(parentRef))
			Expect(parentStatus.ControllerName).To(Equal(gatewayv1.GatewayController(controllerName)))

			accepted := meta.FindStatusCondition(parentStatus.Conditions, string(gatewayv1.RouteConditionAccepted))
			Expect(accepted).NotTo(BeNil())
			Expect(accepted.Status).To(Equal(metav1.ConditionTrue))
			Expect(accepted.Reason).To(Equal(string(gatewayv1.RouteReasonAccepted)))
			Expect(accepted.ObservedGeneration).To(Equal(int64(3)))
		})

		It("updates the existing parent status for this controller", func() {
			group := gatewayv1.Group("gateway.networking.k8s.io")
			kind := gatewayv1.Kind("Gateway")
			namespace := gatewayv1.Namespace("default")
			parentRef := gatewayv1.ParentReference{
				Namespace: &namespace,
				Name:      gatewayv1.ObjectName("gateway"),
			}
			defaultedParentRef := parentRef
			defaultedParentRef.Group = &group
			defaultedParentRef.Kind = &kind
			route := &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Generation: 4,
				},
				Status: gatewayv1.HTTPRouteStatus{
					RouteStatus: gatewayv1.RouteStatus{
						Parents: []gatewayv1.RouteParentStatus{{
							ParentRef:      defaultedParentRef,
							ControllerName: gatewayv1.GatewayController(controllerName),
							Conditions: []metav1.Condition{{
								Type:               string(gatewayv1.RouteConditionAccepted),
								Status:             metav1.ConditionUnknown,
								Reason:             string(gatewayv1.RouteReasonPending),
								ObservedGeneration: 3,
								Message:            "Reconciling with Cloudflare API",
							}},
						}},
					},
				},
			}

			setHTTPRouteParentStatusCondition(route, parentRef, metav1.Condition{
				Type:               string(gatewayv1.RouteConditionAccepted),
				Status:             metav1.ConditionTrue,
				Reason:             string(gatewayv1.RouteReasonAccepted),
				ObservedGeneration: route.Generation,
				Message:            "Successfully reconciled with Cloudflare",
			})

			Expect(route.Status.Parents).To(HaveLen(1))
			accepted := meta.FindStatusCondition(route.Status.Parents[0].Conditions, string(gatewayv1.RouteConditionAccepted))
			Expect(accepted).NotTo(BeNil())
			Expect(accepted.Status).To(Equal(metav1.ConditionTrue))
			Expect(accepted.Reason).To(Equal(string(gatewayv1.RouteReasonAccepted)))
			Expect(accepted.ObservedGeneration).To(Equal(int64(4)))
			Expect(accepted.Message).To(Equal("Successfully reconciled with Cloudflare"))
		})
	})
})
