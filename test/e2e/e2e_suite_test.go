package e2e

import (
	"flag"
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/gateway-api/conformance"
)

// Run e2e tests using the Ginkgo runner.
func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	fmt.Fprintf(GinkgoWriter, "Starting cloudflare-kubernetes-gateway suite\n") //nolint:errcheck
	RunSpecs(t, "e2e suite")

	fmt.Fprintf(GinkgoWriter, "Starting gateway-api conformance suite\n") //nolint:errcheck
	if flag.Set("cleanup-base-resources", "false") != nil || flag.Set("conformance-profiles", "GATEWAY-HTTP") != nil {
		t.Fatal("Failed to set conformance flags")
	}
	conformance.RunConformance(t)
}
