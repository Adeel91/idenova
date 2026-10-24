package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"idenova/document/internal/application"
	"idenova/document/internal/domain"
)

type fakeDocumentRepository struct {
	documents map[string]*domain.Document
}

func (r *fakeDocumentRepository) Save(
	_ context.Context,
	document *domain.Document,
) error {
	if _, exists := r.documents[document.ID]; exists {
		return domain.ErrDocumentDuplicate
	}

	copy := *document
	r.documents[document.ID] = &copy
	return nil
}

func (r *fakeDocumentRepository) FindByID(
	_ context.Context,
	id string,
) (*domain.Document, error) {
	document, exists := r.documents[id]
	if !exists {
		return nil, domain.ErrDocumentNotFound
	}

	copy := *document
	return &copy, nil
}

func (r *fakeDocumentRepository) ListByVerification(
	_ context.Context,
	verificationID string,
) ([]domain.Document, error) {
	results := make([]domain.Document, 0)

	for _, document := range r.documents {
		if document.VerificationID == verificationID {
			results = append(results, *document)
		}
	}

	return results, nil
}

func (r *fakeDocumentRepository) Update(
	_ context.Context,
	document *domain.Document,
) error {
	existing, exists := r.documents[document.ID]
	if !exists {
		return domain.ErrDocumentNotFound
	}

	if existing.Status != domain.StatusPending {
		return domain.ErrInvalidTransition
	}

	copy := *document
	r.documents[document.ID] = &copy
	return nil
}

func TestDocumentHandler(t *testing.T) {
	repo := &fakeDocumentRepository{
		documents: make(map[string]*domain.Document),
	}

	handler := NewHandler(
		application.NewCreateDocument(repo),
		application.NewGetDocument(repo),
		application.NewListDocuments(repo),
		application.NewUpdateDocument(repo),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /documents", handler.CreateDocument)
	mux.HandleFunc("GET /documents/{id}", handler.GetDocument)
	mux.HandleFunc("GET /verifications/{id}/documents", handler.ListDocuments)
	mux.HandleFunc("PATCH /documents/{id}/status", handler.UpdateDocument)

	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		expectedStatus int
	}{
		{
			name:           "create document",
			method:         http.MethodPost,
			path:           "/documents",
			body:           `{"id":"doc_123","verificationId":"verification_123","type":"SYNTHETIC_ID"}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "get document",
			method:         http.MethodGet,
			path:           "/documents/doc_123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "list documents",
			method:         http.MethodGet,
			path:           "/verifications/verification_123/documents",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "approve document",
			method:         http.MethodPatch,
			path:           "/documents/doc_123/status",
			body:           `{"status":"APPROVED"}`,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "approve again",
			method:         http.MethodPatch,
			path:           "/documents/doc_123/status",
			body:           `{"status":"APPROVED"}`,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "duplicate document",
			method:         http.MethodPost,
			path:           "/documents",
			body:           `{"id":"doc_123","verificationId":"verification_123","type":"SYNTHETIC_ID"}`,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "invalid document",
			method:         http.MethodPost,
			path:           "/documents",
			body:           `{"id":"","verificationId":"verification_123","type":"SYNTHETIC_ID"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid JSON",
			method:         http.MethodPost,
			path:           "/documents",
			body:           `{invalid}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing document",
			method:         http.MethodGet,
			path:           "/documents/missing",
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
					t.Fatal("expected error message")
				}
			}
		})
	}
}
