package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"idenova/verification/internal/adapters/memory"
	"idenova/verification/internal/application"
	"idenova/verification/internal/ports"
)

type fakeIdempotencyStore struct {
	key        string
	hash       string
	customerID string
	result     ports.IdempotencyResult
	err        error
}

func (f *fakeIdempotencyStore) CreateVerificationIdempotent(
	_ context.Context,
	key, hash, customerID string,
) (ports.IdempotencyResult, error) {
	f.key = key
	f.hash = hash
	f.customerID = customerID

	return f.result, f.err
}

func TestCreateVerificationWithIdempotencyKey(t *testing.T) {
	tests := []struct {
		name           string
		result         ports.IdempotencyResult
		repoErr        error
		expectedStatus int
		expectedReplay string
	}{
		{
			name: "first request",
			result: ports.IdempotencyResult{
				ResourceID: "ver_123",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "replayed request",
			result: ports.IdempotencyResult{
				ResourceID: "ver_123",
				Replayed:   true,
			},
			expectedStatus: http.StatusOK,
			expectedReplay: "true",
		},
		{
			name:           "conflicting request",
			repoErr:        ports.ErrIdempotencyConflict,
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := memory.NewRepository()

			handler := NewHandler(
				application.NewCreateVerification(repo),
				application.NewGetVerification(repo),
				application.NewListVerifications(repo),
				application.NewUpdateVerification(repo, testDocumentReader{}),
			)

			store := &fakeIdempotencyStore{
				result: tt.result,
				err:    tt.repoErr,
			}

			handler.SetIdempotentCreation(
				application.NewCreateVerificationIdempotent(store),
			)

			mux := http.NewServeMux()
			mux.HandleFunc("POST /verifications", handler.CreateVerification)

			req := httptest.NewRequest(
				http.MethodPost,
				"/verifications",
				strings.NewReader(`{"customerId":"customer_123"}`),
			)
			req.Header.Set("Idempotency-Key", "key_123")

			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf(
					"expected HTTP %d, got %d: %s",
					tt.expectedStatus,
					recorder.Code,
					recorder.Body.String(),
				)
			}

			if got := recorder.Header().Get("Idempotency-Replayed"); got != tt.expectedReplay {
				t.Errorf("expected replay header %q, got %q", tt.expectedReplay, got)
			}

			if tt.repoErr == nil {
				var body struct {
					ID string `json:"id"`
				}

				if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
					t.Fatal(err)
				}

				if body.ID != "ver_123" {
					t.Errorf("expected ver_123, got %s", body.ID)
				}

				if store.key != "key_123" || store.customerID != "customer_123" || store.hash == "" {
					t.Fatal("incorrect idempotency repository arguments")
				}
			}
		})
	}
}

func TestIdempotencyKeyValidation(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		expectedStatus int
	}{
		{
			name:           "empty key uses legacy creation",
			key:            "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "whitespace key",
			key:            "   ",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "key exceeds maximum length",
			key:            strings.Repeat("x", 256),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := memory.NewRepository()

			handler := NewHandler(
				application.NewCreateVerification(repo),
				application.NewGetVerification(repo),
				application.NewListVerifications(repo),
				application.NewUpdateVerification(repo, testDocumentReader{}),
			)

			handler.SetIdempotentCreation(
				application.NewCreateVerificationIdempotent(
					&fakeIdempotencyStore{},
				),
			)

			mux := http.NewServeMux()
			mux.HandleFunc("POST /verifications", handler.CreateVerification)

			req := httptest.NewRequest(
				http.MethodPost,
				"/verifications",
				strings.NewReader(`{"customerId":"customer_123"}`),
			)

			if tt.key != "" {
				req.Header.Set("Idempotency-Key", tt.key)
			}

			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf(
					"expected HTTP %d, got %d: %s",
					tt.expectedStatus,
					recorder.Code,
					recorder.Body.String(),
				)
			}
		})
	}
}
