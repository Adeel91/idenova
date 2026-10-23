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
	"idenova/verification/internal/domain"
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

func TestGetVerification(t *testing.T) {
	repo := memory.NewRepository()

	handler := NewHandler(
		application.NewCreateVerification(repo),
		application.NewGetVerification(repo),
		application.NewListVerifications(repo),
		application.NewUpdateVerification(repo),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /verifications/{id}", handler.GetVerification)

	ctx := context.Background()

	err := repo.Save(ctx, &domain.Verification{
		ID:         "123",
		CustomerID: "456",
		Status:     domain.StatusPending,
	})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		id             string
		expectedStatus int
	}{
		{"existing verification", "123", http.StatusOK},
		{"missing verification", "999", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				"/verifications/"+tt.id,
				nil,
			)

			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Errorf(
					"expected HTTP %d, got %d",
					tt.expectedStatus,
					recorder.Code,
				)
			}
		})
	}
}

func TestUpdateVerification(t *testing.T) {
	repo := memory.NewRepository()

	handler := NewHandler(
		application.NewCreateVerification(repo),
		application.NewGetVerification(repo),
		application.NewListVerifications(repo),
		application.NewUpdateVerification(repo),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("PATCH /verifications/{id}/status", handler.UpdateVerification)

	if err := repo.Save(context.Background(), &domain.Verification{
		ID:         "123",
		CustomerID: "456",
		Status:     domain.StatusPending,
	}); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		id             string
		body           string
		expectedStatus int
	}{
		{"approve pending", "123", `{"status":"APPROVED"}`, http.StatusNoContent},
		{"approve again", "123", `{"status":"APPROVED"}`, http.StatusConflict},
		{"missing verification", "999", `{"status":"APPROVED"}`, http.StatusNotFound},
		{"invalid status", "123", `{"status":"UNKNOWN"}`, http.StatusBadRequest},
		{"invalid JSON", "123", `{invalid}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPatch,
				"/verifications/"+tt.id+"/status",
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

func TestListVerifications(t *testing.T) {
	repo := memory.NewRepository()
	ctx := context.Background()

	records := []domain.Verification{
		{ID: "001", CustomerID: "alice", Status: domain.StatusPending},
		{ID: "002", CustomerID: "bob", Status: domain.StatusApproved},
		{ID: "003", CustomerID: "charlie", Status: domain.StatusPending},
	}

	for i := range records {
		if err := repo.Save(ctx, &records[i]); err != nil {
			t.Fatal(err)
		}
	}

	handler := NewHandler(
		application.NewCreateVerification(repo),
		application.NewGetVerification(repo),
		application.NewListVerifications(repo),
		application.NewUpdateVerification(repo),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /verifications", handler.ListVerifications)

	tests := []struct {
		name           string
		query          string
		expectedStatus int
		expectedCount  int
	}{
		{"all records", "", http.StatusOK, 3},
		{"pending only", "?status=PENDING", http.StatusOK, 2},
		{"pagination", "?limit=1&offset=0", http.StatusOK, 1},
		{"invalid limit", "?limit=200", http.StatusBadRequest, 0},
		{"invalid offset", "?offset=abc", http.StatusBadRequest, 0},
		{"invalid status", "?status=UNKNOWN", http.StatusBadRequest, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				"/verifications"+tt.query,
				nil,
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

			if tt.expectedStatus == http.StatusOK {
				var results []domain.Verification

				if err := json.Unmarshal(recorder.Body.Bytes(), &results); err != nil {
					t.Fatal(err)
				}

				if len(results) != tt.expectedCount {
					t.Errorf(
						"expected %d records, got %d",
						tt.expectedCount,
						len(results),
					)
				}
			}
		})
	}
}
