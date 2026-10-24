package documenthttp

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"idenova/verification/internal/ports"
)

func TestListByVerificationID(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantCount  int
		wantStatus ports.DocumentStatus
		wantError  bool
	}{
		{
			name:       "approved document",
			statusCode: http.StatusOK,
			body:       `[{"ID":"doc_1","Status":"APPROVED"}]`,
			wantCount:  1,
			wantStatus: ports.DocumentApproved,
		},
		{
			name:       "multiple documents",
			statusCode: http.StatusOK,
			body:       `[{"ID":"doc_1","Status":"PENDING"},{"ID":"doc_2","Status":"REJECTED"}]`,
			wantCount:  2,
			wantStatus: ports.DocumentPending,
		},
		{
			name:       "empty documents",
			statusCode: http.StatusOK,
			body:       `[]`,
			wantCount:  0,
		},
		{
			name:       "document service unavailable",
			statusCode: http.StatusServiceUnavailable,
			body:       `{"error":"Service unavailable"}`,
			wantError:  true,
		},
		{
			name:       "invalid JSON",
			statusCode: http.StatusOK,
			body:       `{invalid}`,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					if r.Method != http.MethodGet {
						t.Errorf("expected GET, got %s", r.Method)
					}

					if r.URL.Path != "/verifications/ver_123/documents" {
						t.Errorf("unexpected request path: %s", r.URL.Path)
					}

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(tt.statusCode)

					_, _ = fmt.Fprint(w, tt.body)
				},
			))
			defer server.Close()

			reader := NewReader(server.URL)

			documents, err := reader.ListByVerificationID(
				context.Background(),
				"ver_123",
			)

			if tt.wantError {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(documents) != tt.wantCount {
				t.Fatalf(
					"expected %d documents, got %d",
					tt.wantCount,
					len(documents),
				)
			}

			if tt.wantCount > 0 &&
				documents[0].Status != tt.wantStatus {
				t.Errorf(
					"expected %s, got %s",
					tt.wantStatus,
					documents[0].Status,
				)
			}
		})
	}
}
