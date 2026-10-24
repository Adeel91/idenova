package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"idenova/identity/internal/application"
	"idenova/identity/internal/domain"
)

type fakeOrganizationRepository struct {
	organizations map[string]*domain.Organization
	saveErr       error
	findErr       error
}

func (r *fakeOrganizationRepository) Save(
	_ context.Context,
	org *domain.Organization,
) error {
	if r.saveErr != nil {
		return r.saveErr
	}

	if _, exists := r.organizations[org.ID]; exists {
		return domain.ErrOrganizationDuplicate
	}

	r.organizations[org.ID] = org
	return nil
}

func (r *fakeOrganizationRepository) FindByID(
	_ context.Context,
	id string,
) (*domain.Organization, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}

	org, exists := r.organizations[id]
	if !exists {
		return nil, domain.ErrOrganizationNotFound
	}

	return org, nil
}

func TestOrganizationHandler(t *testing.T) {
	repo := &fakeOrganizationRepository{
		organizations: make(map[string]*domain.Organization),
	}

	handler := NewHandler(
		application.NewCreateOrganization(repo),
		application.NewGetOrganization(repo),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /organizations", handler.CreateOrganization)
	mux.HandleFunc("GET /organizations/{id}", handler.GetOrganization)

	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		expectedStatus int
	}{
		{
			name:           "create organization",
			method:         http.MethodPost,
			path:           "/organizations",
			body:           `{"id":"org_123","name":"Acme GmbH"}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "get organization",
			method:         http.MethodGet,
			path:           "/organizations/org_123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "duplicate organization",
			method:         http.MethodPost,
			path:           "/organizations",
			body:           `{"id":"org_123","name":"Acme GmbH"}`,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "invalid organization",
			method:         http.MethodPost,
			path:           "/organizations",
			body:           `{"id":"","name":"Acme GmbH"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid JSON",
			method:         http.MethodPost,
			path:           "/organizations",
			body:           `{invalid}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "organization not found",
			method:         http.MethodGet,
			path:           "/organizations/missing",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				tt.method,
				tt.path,
				strings.NewReader(tt.body),
			)

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

			if recorder.Code >= 400 {
				var response map[string]string

				if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
					t.Fatalf("invalid JSON error response: %v", err)
				}

				if response["error"] == "" {
					t.Fatal("expected a nonempty error message")
				}
			}
		})
	}
}

func TestOrganizationHandlerInternalErrors(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		path    string
		body    string
		saveErr error
		findErr error
	}{
		{
			name:    "create repository failure",
			method:  http.MethodPost,
			path:    "/organizations",
			body:    `{"id":"org_123","name":"Acme GmbH"}`,
			saveErr: errors.New("database unavailable"),
		},
		{
			name:    "get repository failure",
			method:  http.MethodGet,
			path:    "/organizations/org_123",
			findErr: errors.New("database unavailable"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeOrganizationRepository{
				organizations: make(map[string]*domain.Organization),
				saveErr:       tt.saveErr,
				findErr:       tt.findErr,
			}

			handler := NewHandler(
				application.NewCreateOrganization(repo),
				application.NewGetOrganization(repo),
			)

			mux := http.NewServeMux()
			mux.HandleFunc("POST /organizations", handler.CreateOrganization)
			mux.HandleFunc("GET /organizations/{id}", handler.GetOrganization)

			req := httptest.NewRequest(
				tt.method,
				tt.path,
				strings.NewReader(tt.body),
			)
			recorder := httptest.NewRecorder()

			mux.ServeHTTP(recorder, req)

			if recorder.Code != http.StatusInternalServerError {
				t.Fatalf(
					"expected HTTP 500, got %d: %s",
					recorder.Code,
					recorder.Body.String(),
				)
			}

			var response map[string]string
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Fatal(err)
			}

			if response["error"] != "Internal server error" {
				t.Fatalf("unexpected error response: %v", response)
			}
		})
	}
}
