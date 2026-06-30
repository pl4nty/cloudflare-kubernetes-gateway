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
			parentRef := gatewayv1.ParentReference{
				Name: gatewayv1.ObjectName("gateway"),
			}
			otherParentRef := gatewayv1.ParentReference{
				Name: gatewayv1.ObjectName("other-gateway"),
			}
			route := &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Generation: 3,
					Namespace:  "default",
				},
				Status: gatewayv1.HTTPRouteStatus{
					RouteStatus: gatewayv1.RouteStatus{
						Parents: []gatewayv1.RouteParentStatus{{
							ParentRef:      normalizeHTTPRouteParentRef(otherParentRef, "default"),
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
			Expect(findHTTPRouteParentStatus(route, gatewayv1.GatewayController("example.com/other-controller"), otherParentRef)).NotTo(BeNil())

			parentStatus := findHTTPRouteParentStatus(route, gatewayv1.GatewayController(controllerName), parentRef)
			Expect(parentStatus).NotTo(BeNil())
			Expect(parentStatus.ParentRef).To(Equal(normalizeHTTPRouteParentRef(parentRef, route.Namespace)))
			Expect(parentStatus.ControllerName).To(Equal(gatewayv1.GatewayController(controllerName)))

			accepted := meta.FindStatusCondition(parentStatus.Conditions, string(gatewayv1.RouteConditionAccepted))
			Expect(accepted).NotTo(BeNil())
			Expect(accepted.Status).To(Equal(metav1.ConditionTrue))
			Expect(accepted.Reason).To(Equal(string(gatewayv1.RouteReasonAccepted)))
			Expect(accepted.ObservedGeneration).To(Equal(int64(3)))
		})

		It("updates the existing parent status for this controller", func() {
			parentRef := gatewayv1.ParentReference{
				Name: gatewayv1.ObjectName("gateway"),
			}
			defaultedParentRef := normalizeHTTPRouteParentRef(parentRef, "default")
			route := &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Generation: 4,
					Namespace:  "default",
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
			parentStatus := findHTTPRouteParentStatus(route, gatewayv1.GatewayController(controllerName), parentRef)
			Expect(parentStatus).NotTo(BeNil())
			accepted := meta.FindStatusCondition(parentStatus.Conditions, string(gatewayv1.RouteConditionAccepted))
			Expect(accepted).NotTo(BeNil())
			Expect(accepted.Status).To(Equal(metav1.ConditionTrue))
			Expect(accepted.Reason).To(Equal(string(gatewayv1.RouteReasonAccepted)))
			Expect(accepted.ObservedGeneration).To(Equal(int64(4)))
			Expect(accepted.Message).To(Equal("Successfully reconciled with Cloudflare"))
		})

		It("sets Accepted false for unsupported parentRef groups", func() {
			group := gatewayv1.Group("core")
			parentRef := gatewayv1.ParentReference{
				Group: &group,
				Name:  gatewayv1.ObjectName("gateway"),
			}
			err := validateHTTPRouteGatewayParentRef(parentRef)
			Expect(err).To(MatchError("unsupported parentRef group core"))

			route := &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Generation: 5,
					Namespace:  "default",
				},
			}
			setHTTPRouteParentStatusCondition(route, parentRef, metav1.Condition{
				Type:               string(gatewayv1.RouteConditionAccepted),
				Status:             metav1.ConditionFalse,
				Reason:             string(gatewayv1.RouteReasonNoMatchingParent),
				ObservedGeneration: route.Generation,
				Message:            "ParentRef validation failed: unsupported parentRef group core",
			})

			parentStatus := findHTTPRouteParentStatus(route, gatewayv1.GatewayController(controllerName), parentRef)
			Expect(parentStatus).NotTo(BeNil())
			accepted := meta.FindStatusCondition(parentStatus.Conditions, string(gatewayv1.RouteConditionAccepted))
			Expect(accepted).NotTo(BeNil())
			Expect(accepted.Status).To(Equal(metav1.ConditionFalse))
			Expect(accepted.Reason).To(Equal(string(gatewayv1.RouteReasonNoMatchingParent)))
			Expect(accepted.Message).To(Equal("ParentRef validation failed: unsupported parentRef group core"))
		})

		It("sets Accepted false for unsupported parentRef kinds", func() {
			kind := gatewayv1.Kind("Service")
			parentRef := gatewayv1.ParentReference{
				Kind: &kind,
				Name: gatewayv1.ObjectName("service"),
			}
			err := validateHTTPRouteGatewayParentRef(parentRef)
			Expect(err).To(MatchError("unsupported parentRef kind Service"))

			route := &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Generation: 6,
					Namespace:  "default",
				},
			}
			setHTTPRouteParentStatusCondition(route, parentRef, metav1.Condition{
				Type:               string(gatewayv1.RouteConditionAccepted),
				Status:             metav1.ConditionFalse,
				Reason:             string(gatewayv1.RouteReasonNoMatchingParent),
				ObservedGeneration: route.Generation,
				Message:            "ParentRef validation failed: unsupported parentRef kind Service",
			})

			parentStatus := findHTTPRouteParentStatus(route, gatewayv1.GatewayController(controllerName), parentRef)
			Expect(parentStatus).NotTo(BeNil())
			accepted := meta.FindStatusCondition(parentStatus.Conditions, string(gatewayv1.RouteConditionAccepted))
			Expect(accepted).NotTo(BeNil())
			Expect(accepted.Status).To(Equal(metav1.ConditionFalse))
			Expect(accepted.Reason).To(Equal(string(gatewayv1.RouteReasonNoMatchingParent)))
			Expect(accepted.Message).To(Equal("ParentRef validation failed: unsupported parentRef kind Service"))
		})
	})
})

func findHTTPRouteParentStatus(route *gatewayv1.HTTPRoute, controller gatewayv1.GatewayController, parentRef gatewayv1.ParentReference) *gatewayv1.RouteParentStatus {
	normalizedParentRef := normalizeHTTPRouteParentRef(parentRef, route.Namespace)
	for i := range route.Status.Parents {
		parentStatus := &route.Status.Parents[i]
		if parentStatus.ControllerName == controller && parentRefsEqual(normalizeHTTPRouteParentRef(parentStatus.ParentRef, route.Namespace), normalizedParentRef) {
			return parentStatus
		}
	}
	return nil
}
