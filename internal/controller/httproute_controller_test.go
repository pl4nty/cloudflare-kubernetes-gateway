// SPDX-FileCopyrightText: 2023-2026 Tom Plant, Elias Elwyn, and contributors
// SPDX-License-Identifier: MIT
package controller

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

var _ = Describe("HTTPRoute Controller", func() {
	Context("When reconciling a resource", func() {

		It("should successfully reconcile the resource", func() {

			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})

	// updateRouteStatus doesn't call the Cloudflare API, so it can be
	// exercised directly against envtest without credentials.
	Context("updateRouteStatus", func() {
		var (
			ctx        context.Context
			reconciler *HTTPRouteReconciler
			gwClass    *gatewayv1.GatewayClass
			gw         *gatewayv1.Gateway
			svc        *corev1.Service
		)

		BeforeEach(func() {
			ctx = context.Background()
			reconciler = &HTTPRouteReconciler{Client: k8sClient, Scheme: k8sClient.Scheme()}

			gwClass = &gatewayv1.GatewayClass{
				ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("cloudflare-%d", GinkgoParallelProcess())},
				Spec: gatewayv1.GatewayClassSpec{
					ControllerName: controllerName,
				},
			}
			Expect(k8sClient.Create(ctx, gwClass)).To(Succeed())

			gw = &gatewayv1.Gateway{
				ObjectMeta: metav1.ObjectMeta{Name: "gateway", Namespace: "default"},
				Spec: gatewayv1.GatewaySpec{
					GatewayClassName: gatewayv1.ObjectName(gwClass.Name),
					Listeners: []gatewayv1.Listener{
						{Name: "http", Protocol: gatewayv1.HTTPProtocolType, Port: 80},
					},
				},
			}
			Expect(k8sClient.Create(ctx, gw)).To(Succeed())

			svc = &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{Name: "whoami", Namespace: "default"},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{{Port: 80}},
				},
			}
			Expect(k8sClient.Create(ctx, svc)).To(Succeed())
		})

		AfterEach(func() {
			_ = k8sClient.Delete(ctx, gw)
			_ = k8sClient.Delete(ctx, gwClass)
			_ = k8sClient.Delete(ctx, svc)
		})

		It("sets Accepted/ResolvedRefs=True when the backend Service exists", func() {
			route := &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{Name: "whoami", Namespace: "default"},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{{Name: gatewayv1.ObjectName(gw.Name)}},
					},
					Hostnames: []gatewayv1.Hostname{"whoami.example.com"},
					Rules: []gatewayv1.HTTPRouteRule{{
						BackendRefs: []gatewayv1.HTTPBackendRef{{
							BackendRef: gatewayv1.BackendRef{
								BackendObjectReference: gatewayv1.BackendObjectReference{
									Name: gatewayv1.ObjectName(svc.Name),
									Port: ptrPort(80),
								},
							},
						}},
					}},
				},
			}
			Expect(k8sClient.Create(ctx, route)).To(Succeed())
			defer func() { _ = k8sClient.Delete(ctx, route) }()

			Expect(reconciler.updateRouteStatus(ctx, route)).To(Succeed())

			got := &gatewayv1.HTTPRoute{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: route.Name, Namespace: route.Namespace}, got)).To(Succeed())

			Expect(got.Status.Parents).To(HaveLen(1))
			parent := got.Status.Parents[0]
			Expect(string(parent.ControllerName)).To(Equal(controllerName))
			Expect(parent.ParentRef.Name).To(Equal(gatewayv1.ObjectName(gw.Name)))

			acceptedCond := findCondition(parent.Conditions, string(gatewayv1.RouteConditionAccepted))
			Expect(acceptedCond).NotTo(BeNil())
			Expect(acceptedCond.Status).To(Equal(metav1.ConditionTrue))

			resolvedCond := findCondition(parent.Conditions, string(gatewayv1.RouteConditionResolvedRefs))
			Expect(resolvedCond).NotTo(BeNil())
			Expect(resolvedCond.Status).To(Equal(metav1.ConditionTrue))
		})

		It("sets ResolvedRefs=False with reason BackendNotFound when the backend Service is missing", func() {
			route := &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{Name: "missing-backend", Namespace: "default"},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{{Name: gatewayv1.ObjectName(gw.Name)}},
					},
					Hostnames: []gatewayv1.Hostname{"missing.example.com"},
					Rules: []gatewayv1.HTTPRouteRule{{
						BackendRefs: []gatewayv1.HTTPBackendRef{{
							BackendRef: gatewayv1.BackendRef{
								BackendObjectReference: gatewayv1.BackendObjectReference{
									Name: "does-not-exist",
									Port: ptrPort(80),
								},
							},
						}},
					}},
				},
			}
			Expect(k8sClient.Create(ctx, route)).To(Succeed())
			defer func() { _ = k8sClient.Delete(ctx, route) }()

			Expect(reconciler.updateRouteStatus(ctx, route)).To(Succeed())

			got := &gatewayv1.HTTPRoute{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: route.Name, Namespace: route.Namespace}, got)).To(Succeed())

			Expect(got.Status.Parents).To(HaveLen(1))
			resolvedCond := findCondition(got.Status.Parents[0].Conditions, string(gatewayv1.RouteConditionResolvedRefs))
			Expect(resolvedCond).NotTo(BeNil())
			Expect(resolvedCond.Status).To(Equal(metav1.ConditionFalse))
			Expect(resolvedCond.Reason).To(Equal(string(gatewayv1.RouteReasonBackendNotFound)))
		})

		It("leaves status.parents entries owned by other controllers untouched", func() {
			route := &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{Name: "multi-parent", Namespace: "default"},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{{Name: gatewayv1.ObjectName(gw.Name)}},
					},
					Hostnames: []gatewayv1.Hostname{"multi.example.com"},
					Rules: []gatewayv1.HTTPRouteRule{{
						BackendRefs: []gatewayv1.HTTPBackendRef{{
							BackendRef: gatewayv1.BackendRef{
								BackendObjectReference: gatewayv1.BackendObjectReference{
									Name: gatewayv1.ObjectName(svc.Name),
									Port: ptrPort(80),
								},
							},
						}},
					}},
				},
			}
			Expect(k8sClient.Create(ctx, route)).To(Succeed())
			defer func() { _ = k8sClient.Delete(ctx, route) }()

			route.Status.Parents = []gatewayv1.RouteParentStatus{{
				ParentRef:      gatewayv1.ParentReference{Name: "some-other-gateway"},
				ControllerName: "example.com/other-controller",
				Conditions: []metav1.Condition{{
					Type:               string(gatewayv1.RouteConditionAccepted),
					Status:             metav1.ConditionTrue,
					Reason:             string(gatewayv1.RouteReasonAccepted),
					ObservedGeneration: route.Generation,
					LastTransitionTime: metav1.Now(),
				}},
			}}
			Expect(k8sClient.Status().Update(ctx, route)).To(Succeed())

			Expect(reconciler.updateRouteStatus(ctx, route)).To(Succeed())

			got := &gatewayv1.HTTPRoute{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: route.Name, Namespace: route.Namespace}, got)).To(Succeed())
			Expect(got.Status.Parents).To(HaveLen(2))

			var sawOther, sawOurs bool
			for _, p := range got.Status.Parents {
				if string(p.ControllerName) == "example.com/other-controller" {
					sawOther = true
				}
				if string(p.ControllerName) == controllerName {
					sawOurs = true
				}
			}
			Expect(sawOther).To(BeTrue())
			Expect(sawOurs).To(BeTrue())
		})
	})
})

func findCondition(conditions []metav1.Condition, condType string) *metav1.Condition {
	for i := range conditions {
		if conditions[i].Type == condType {
			return &conditions[i]
		}
	}
	return nil
}

func ptrPort(p gatewayv1.PortNumber) *gatewayv1.PortNumber {
	return &p
}
