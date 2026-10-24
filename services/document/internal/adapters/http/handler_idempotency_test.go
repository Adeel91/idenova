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
	"idenova/idempotency"
)

type handlerIdempotencyRepository struct {
	documents *fakeDocumentRepository
	requests  map[string]struct {
		hash string
		id   string
	}
}

func (f *handlerIdempotencyRepository) CreateDocumentIdempotent(
	ctx context.Context,
	key, hash, verificationID, documentType string,
) (idempotency.Result, error) {
	if existing, ok := f.requests[key]; ok {
		if existing.hash != hash {
			return idempotency.Result{}, idempotency.ErrConflict
		}

		return idempotency.Result{
			ResourceID: existing.id,
			Replayed:   true,
		}, nil
	}

	id := "doc_test_" + key

	document, err := domain.NewDocument(id, verificationID, documentType)
	if err != nil {
		return idempotency.Result{}, err
	}

	if err := f.documents.Save(ctx, document); err != nil {
		return idempotency.Result{}, err
	}

	f.requests[key] = struct {
		hash string
		id   string
	}{hash: hash, id: id}

	return idempotency.Result{ResourceID: id}, nil
}

func TestIdempotentDocumentHTTPHandler(t *testing.T) {
	repo := &fakeDocumentRepository{
		documents: make(map[string]*domain.Document),
	}

	idemRepo := &handlerIdempotencyRepository{
		documents: repo,
		requests: make(map[string]struct {
			hash string
			id   string
		}),
	}

	handler := NewHandler(
		application.NewCreateDocument(repo),
		application.NewGetDocument(repo),
		application.NewListDocuments(repo),
		application.NewUpdateDocument(repo),
	)

	handler.SetIdempotentCreation(
		application.NewCreateDocumentIdempotent(idemRepo),
	)

	send := func(key, body string) *httptest.ResponseRecorder {
		t.Helper()

		req := httptest.NewRequest(
			http.MethodPost,
			"/documents",
			strings.NewReader(body),
		)
		req.Header.Set("Idempotency-Key", key)

		rec := httptest.NewRecorder()
		handler.CreateDocument(rec, req)
		return rec
	}

	const original = `{"verificationId":"ver_123","type":"SYNTHETIC_ID"}`

	first := send("key_123", original)
	if first.Code != http.StatusCreated {
		t.Fatalf("first request: expected 201, got %d: %s",
			first.Code, first.Body.String())
	}

	var created domain.Document
	if err := json.Unmarshal(first.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if created.ID == "" {
		t.Fatal("first response missing document ID")
	}

	replay := send("key_123", original)
	if replay.Code != http.StatusOK {
		t.Fatalf("replay: expected 200, got %d: %s",
			replay.Code, replay.Body.String())
	}
	if replay.Header().Get("Idempotency-Replayed") != "true" {
		t.Fatal("missing replay response header")
	}

	var repeated domain.Document
	if err := json.Unmarshal(replay.Body.Bytes(), &repeated); err != nil {
		t.Fatal(err)
	}
	if repeated.ID != created.ID {
		t.Fatal("replay returned a different document ID")
	}

	typeConflict := send(
		"key_123",
		`{"verificationId":"ver_123","type":"SYNTHETIC_PASSPORT"}`,
	)
	if typeConflict.Code != http.StatusConflict {
		t.Fatalf("type conflict: expected 409, got %d",
			typeConflict.Code)
	}

	verificationConflict := send(
		"key_123",
		`{"verificationId":"ver_456","type":"SYNTHETIC_ID"}`,
	)
	if verificationConflict.Code != http.StatusConflict {
		t.Fatalf("verification conflict: expected 409, got %d",
			verificationConflict.Code)
	}

	invalidKey := send("bad key", original)
	if invalidKey.Code != http.StatusBadRequest {
		t.Fatalf("invalid key: expected 400, got %d",
			invalidKey.Code)
	}

	invalidDocument := send(
		"key_invalid_document",
		`{"verificationId":"","type":"SYNTHETIC_ID"}`,
	)
	if invalidDocument.Code != http.StatusBadRequest {
		t.Fatalf("invalid document: expected 400, got %d",
			invalidDocument.Code)
	}

	if len(repo.documents) != 1 {
		t.Fatalf("expected exactly one document, got %d",
			len(repo.documents))
	}
}
