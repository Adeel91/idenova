package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"idenova/verification/internal/adapters/memory"
	"idenova/verification/internal/application"
)

func TestCreateVerification(t *testing.T) {
	repo := memory.NewRepository()

	handler := NewHandler(
		application.NewCreateVerification(repo),
		application.NewGetVerification(repo),
		application.NewListVerifications(repo),
		application.NewUpdateVerification(repo),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /verifications", handler.CreateVerification)

	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "valid verification",
			body:           `{"id":"123","customerId":"456"}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "duplicate verification",
			body:           `{"id":"123","customerId":"456"}`,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "missing customer ID",
			body:           `{"id":"789"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid JSON",
			body:           `{invalid}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPost,
				"/verifications",
				strings.NewReader(tt.body),
			)

			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Errorf(
					"expected HTTP %d, got %d: %s",
					tt.expectedStatus,
					recorder.Code,
					recorder.Body.String(),
				)
			}
		})
	}
}
