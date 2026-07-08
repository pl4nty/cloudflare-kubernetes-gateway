package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/pl4nty/cloudflare-kubernetes-gateway/test/utils"
)

const namespace = "cloudflare-gateway"

var _ = Describe("controller", Ordered, func() {
	BeforeAll(func() {
		By("installing prometheus operator")
		Expect(utils.InstallPrometheusOperator()).To(Succeed())

		By("installing the cert-manager")
		Expect(utils.InstallCertManager()).To(Succeed())

		By("creating manager namespace")
		cmd := exec.Command("kubectl", "create", "ns", namespace)
		_, _ = utils.Run(cmd)
	})

	Context("Operator", func() {
		It("should run successfully", func() {
			var controllerPodName string
			var err error

			imageName, ok := os.LookupEnv("IMAGE_NAME")
			Ω(ok).Should(BeTrueBecause("IMAGE_NAME env var should exist"))
			imageTag, ok :=  os.LookupEnv("IMAGE_TAG")
			Ω(ok).Should(BeTrueBecause("IMAGE_TAG env var should exist"))

			// projectimage stores the name of the image used in the example
			var projectimage = imageName + ":" + imageTag

			By("checking if the manager(Operator) image exists locally")
			cmd := exec.Command("make", "docker-image-ls", fmt.Sprintf("IMG=%s", projectimage))
			_, err = utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			By("loading the the manager(Operator) image on Kind")
			err = utils.LoadImageToKindClusterWithName(projectimage)
			Expect(err).NotTo(HaveOccurred())

			By("installing CRDs")
			cmd = exec.Command("make", "install")
			_, err = utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			By("deploying the controller-manager")
			cmd = exec.Command("make", "deploy", fmt.Sprintf("IMG=%s", projectimage))
			_, err = utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			By("validating that the controller-manager pod is running as expected")
			verifyControllerUp := func() error {
				GinkgoHelper()

				// Get pod name

				cmd = exec.Command("kubectl", "get",
					"pods", "-l", "control-plane=controller-manager",
					"-o", "go-template={{ range .items }}"+
						"{{ if not .metadata.deletionTimestamp }}"+
						"{{ .metadata.name }}"+
						"{{ \"\\n\" }}{{ end }}{{ end }}",
					"-n", namespace,
				)

				podOutput, err := utils.Run(cmd)
				Expect(err).NotTo(HaveOccurred())
				podNames := utils.GetNonEmptyLines(string(podOutput))
				if len(podNames) != 1 {
					return fmt.Errorf("expect 1 controller pods running, but got %d", len(podNames))
				}
				controllerPodName = podNames[0]
				Expect(controllerPodName).Should(ContainSubstring("controller-manager"))

				// Validate pod status
				cmd = exec.Command("kubectl", "get",
					"pods", controllerPodName, "-o", "jsonpath={.status.phase}",
					"-n", namespace,
				)
				status, err := utils.Run(cmd)
				Expect(err).NotTo(HaveOccurred())
				if string(status) != "Running" {
					return fmt.Errorf("controller pod in %s status", status)
				}
				return nil
			}
			Eventually(verifyControllerUp, time.Minute, time.Second).Should(Succeed())

			By("creating the Gateway Secret")
			cmd = exec.Command("kubectl", "create", "secret", "generic", "gateway-conformance",
				"--from-literal=ACCOUNT_ID="+os.Getenv("CLOUDFLARE_ACCOUNT_ID"),
				"--from-literal=TOKEN="+os.Getenv("CLOUDFLARE_API_TOKEN"),
				"-n", namespace)
			_, err = utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			By("creating the GatewayClass")
			cmd = exec.Command("kubectl", "apply", "-f", "test/e2e/gatewayclass.yaml")
			_, err = utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
