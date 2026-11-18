// SPDX-FileCopyrightText: 2023-2026 Tom Plant, Elias Elwyn, and contributors
// SPDX-License-Identifier: MIT
package controller

import (
	"context"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// The kube-controller-manager adds annotations such as
// deployment.kubernetes.io/revision to a Deployment. reconcileDeployment must
// keep them when it updates the Deployment: removing them bumps
// metadata.generation, which re-triggers the reconcile through the Deployment watch.
func TestReconcileDeploymentKeepsServerAnnotations(t *testing.T) {
	const revision = "deployment.kubernetes.io/revision"
	ctx := context.Background()

	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	if err := gatewayv1.Install(scheme); err != nil {
		t.Fatal(err)
	}

	gateway := &gatewayv1.Gateway{ObjectMeta: metav1.ObjectMeta{Name: "gw", Namespace: "ns", UID: "gw-uid"}}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "gw", Namespace: "ns"},
		Data:       map[string][]byte{"TUNNEL_TOKEN": []byte("token")},
	}
	existing := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "gw",
			Namespace:   "ns",
			Annotations: map[string]string{revision: "2"},
		},
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{Spec: corev1.PodSpec{
				Containers: []corev1.Container{{Name: "gateway", Image: "cloudflared:old"}},
			}},
		},
	}

	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(gateway, secret, existing).Build()
	r := &GatewayReconciler{Client: c, Scheme: scheme}
	t.Setenv("GATEWAY_IMAGE", "cloudflared:new")

	for i := 0; i < 2; i++ {
		if _, err := r.reconcileDeployment(ctx, gateway); err != nil {
			t.Fatalf("reconcile %d: %v", i, err)
		}
	}

	got := &appsv1.Deployment{}
	if err := c.Get(ctx, types.NamespacedName{Name: "gw", Namespace: "ns"}, got); err != nil {
		t.Fatal(err)
	}
	if got.Annotations[revision] != "2" {
		t.Errorf("annotation %s was removed, annotations = %v", revision, got.Annotations)
	}
	if image := got.Spec.Template.Spec.Containers[0].Image; image != "cloudflared:new" {
		t.Errorf("image = %q, want cloudflared:new", image)
	}
	if !metav1.IsControlledBy(got, gateway) {
		t.Errorf("Deployment is not owned by the Gateway, ownerReferences = %v", got.OwnerReferences)
	}
}
