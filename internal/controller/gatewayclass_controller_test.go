package controller

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudflare/cloudflare-go/v2"
	"github.com/cloudflare/cloudflare-go/v2/option"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("GatewayClass Controller", func() {
	Context("When reconciling a resource", func() {

		It("should successfully reconcile the resource", func() {

			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})
})

func TestVerifyAccountToken(t *testing.T) {
	const account = "abc123"

	tests := []struct {
		name       string
		statusCode int
		body       string
		wantStatus string
		wantErr    bool
	}{
		{
			name:       "active account token",
			statusCode: http.StatusOK,
			body:       `{"success":true,"errors":[],"messages":[],"result":{"id":"tok","status":"active"}}`,
			wantStatus: "active",
		},
		{
			name:       "disabled account token",
			statusCode: http.StatusOK,
			body:       `{"success":true,"errors":[],"messages":[],"result":{"id":"tok","status":"disabled"}}`,
			wantStatus: "disabled",
		},
		{
			name:       "invalid token returns error",
			statusCode: http.StatusUnauthorized,
			body:       `{"success":false,"errors":[{"code":1000,"message":"Invalid API Token"}],"result":null}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantPath := fmt.Sprintf("/accounts/%s/tokens/verify", account)
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != wantPath {
					t.Errorf("unexpected path: got %q, want %q", r.URL.Path, wantPath)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			api := cloudflare.NewClient(
				option.WithAPIToken("test-token"),
				option.WithBaseURL(srv.URL),
			)

			status, err := verifyAccountToken(context.Background(), api, account)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil (status %q)", status)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if status != tt.wantStatus {
				t.Errorf("status = %q, want %q", status, tt.wantStatus)
			}
		})
	}
}
