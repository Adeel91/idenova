package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"testing"

	"idenova/verification/internal/ports"
)

type fakeIdempotencyRepository struct {
	key        string
	hash       string
	customerID string
	result     ports.IdempotencyResult
	err        error
	called     bool
}

func (f *fakeIdempotencyRepository) CreateVerificationIdempotent(
	_ context.Context,
	key, hash, customerID string,
) (ports.IdempotencyResult, error) {
	f.called = true
	f.key = key
	f.hash = hash
	f.customerID = customerID

	return f.result, f.err
}

func TestCreateVerificationIdempotent(t *testing.T) {
	repo := &fakeIdempotencyRepository{
		result: ports.IdempotencyResult{
			ResourceID: "ver_123",
			Replayed:   false,
		},
	}

	useCase := NewCreateVerificationIdempotent(repo)

	result, err := useCase.Execute(
		context.Background(),
		"key_123",
		"customer_123",
	)
	if err != nil {
		t.Fatal(err)
	}

	if result.ResourceID != "ver_123" || result.Replayed {
		t.Fatalf("unexpected result: %+v", result)
	}

	expectedHash := sha256.Sum256(
		[]byte("create_verification:customer_123"),
	)

	if repo.hash != hex.EncodeToString(expectedHash[:]) {
		t.Fatalf("unexpected request hash: %s", repo.hash)
	}

	if repo.key != "key_123" || repo.customerID != "customer_123" {
		t.Fatalf("unexpected repository arguments")
	}
}

func TestCreateVerificationIdempotentInvalidInput(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		customerID string
	}{
		{"empty key", "", "customer_123"},
		{"long key", strings.Repeat("x", 256), "customer_123"},
		{"empty customer", "key_123", ""},
		{"long customer", "key_123", strings.Repeat("x", 256)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeIdempotencyRepository{}
			useCase := NewCreateVerificationIdempotent(repo)

			_, err := useCase.Execute(
				context.Background(),
				tt.key,
				tt.customerID,
			)

			if err == nil {
				t.Fatal("expected validation error")
			}

			if repo.called {
				t.Fatal("repository should not be called")
			}
		})
	}
}

func TestCreateVerificationIdempotentPropagatesError(t *testing.T) {
	expected := errors.New("database unavailable")

	repo := &fakeIdempotencyRepository{err: expected}
	useCase := NewCreateVerificationIdempotent(repo)

	_, err := useCase.Execute(
		context.Background(),
		"key_123",
		"customer_123",
	)

	if !errors.Is(err, expected) {
		t.Fatalf("expected database error, got %v", err)
	}
}
