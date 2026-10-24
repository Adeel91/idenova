package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"idenova/document/internal/application"
	"idenova/idempotency"
)

type fakeIdempotencyRepository struct {
	key  string
	hash string
	id   string
}

func (f *fakeIdempotencyRepository) CreateDocumentIdempotent(
	_ context.Context,
	key, hash, _, _ string,
) (idempotency.Result, error) {
	if f.key == "" {
		f.key = key
		f.hash = hash
		f.id = "doc_test_123"

		return idempotency.Result{ResourceID: f.id}, nil
	}

	if key == f.key && hash != f.hash {
		return idempotency.Result{}, idempotency.ErrConflict
	}

	return idempotency.Result{
		ResourceID: f.id,
		Replayed:   true,
	}, nil
}

func TestIdempotencyErrorMapping(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		status int
	}{
		{"invalid key", idempotency.ErrInvalidKey, http.StatusBadRequest},
		{"conflicting request", idempotency.ErrConflict, http.StatusConflict},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			handleDomainError(recorder, tt.err)

			if recorder.Code != tt.status {
				t.Fatalf("expected %d, got %d", tt.status, recorder.Code)
			}
		})
	}
}

func TestIdempotencyApplicationFingerprint(t *testing.T) {
	repo := &fakeIdempotencyRepository{}
	useCase := application.NewCreateDocumentIdempotent(repo)

	first, err := useCase.Execute(
		context.Background(), "key_123", "ver_123", "PASSPORT",
	)
	if err != nil {
		t.Fatal(err)
	}

	replay, err := useCase.Execute(
		context.Background(), "key_123", "ver_123", "PASSPORT",
	)
	if err != nil {
		t.Fatal(err)
	}

	if first.ResourceID != replay.ResourceID || !replay.Replayed {
		t.Fatalf("replay failed: first=%+v replay=%+v", first, replay)
	}

	_, err = useCase.Execute(
		context.Background(), "key_123", "ver_123", "ID_CARD",
	)
	if !errors.Is(err, idempotency.ErrConflict) {
		t.Fatalf("expected document type conflict, got %v", err)
	}

	_, err = useCase.Execute(
		context.Background(), "key_123", "ver_456", "PASSPORT",
	)
	if !errors.Is(err, idempotency.ErrConflict) {
		t.Fatalf("expected verification conflict, got %v", err)
	}
}

func TestIdempotencyKeyValidation(t *testing.T) {
	repo := &fakeIdempotencyRepository{}
	useCase := application.NewCreateDocumentIdempotent(repo)

	_, err := useCase.Execute(
		context.Background(), "bad key", "ver_123", "PASSPORT",
	)

	if !errors.Is(err, idempotency.ErrInvalidKey) {
		t.Fatalf("expected invalid key, got %v", err)
	}
}

// Keep the HTTP test file ready for full handler-level tests
// once the existing in-memory document repository setup is reused.
var _ = bytes.NewBuffer
var _ = json.NewDecoder
