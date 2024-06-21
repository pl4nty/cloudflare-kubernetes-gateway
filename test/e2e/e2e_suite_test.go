package e2e

import (
	"fmt"
	"os"
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
	os.Args = []string{"noop", "--cleanup-base-resources=false", "--conformance-profiles=GATEWAY-HTTP"}
	conformance.RunConformance(t)
}
