package postgres

import (
	"context"
	"errors"
	"testing"

	"idenova/document/internal/domain"
	"idenova/idempotency"

	"github.com/google/uuid"
)

type verificationPolicyFunc func(context.Context, string) error

func (f verificationPolicyFunc) RequirePending(
	ctx context.Context,
	verificationID string,
) error {
	return f(ctx, verificationID)
}

func TestDocumentCreationVerificationPolicy(t *testing.T) {
	tests := []struct {
		name    string
		allowed bool
		wantErr error
	}{
		{"PENDING", true, nil},
		{"APPROVED", false, domain.ErrVerificationNotPending},
		{"REJECTED", false, domain.ErrVerificationNotPending},
		{"MISSING", false, domain.ErrVerificationNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := testIdempotencyStore(t)
			calls := 0

			store.verifier = verificationPolicyFunc(
				func(_ context.Context, _ string) error {
					calls++
					return tt.wantErr
				},
			)

			result, err := store.CreateDocumentIdempotent(
				context.Background(),
				uuid.NewString(),
				"request_hash",
				"ver_test_"+uuid.NewString(),
				"SYNTHETIC_ID",
			)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
			if calls != 1 {
				t.Fatalf("expected one verification lookup, got %d", calls)
			}
			if tt.allowed && (result.ResourceID == "" || result.Replayed) {
				t.Fatalf("unexpected creation result: %+v", result)
			}
			if !tt.allowed && result.ResourceID != "" {
				t.Fatal("blocked request created a document")
			}
		})
	}
}

func TestIdempotencyReplayAfterVerificationApproval(t *testing.T) {
	store := testIdempotencyStore(t)
	ctx := context.Background()

	status := "PENDING"
	calls := 0

	store.verifier = verificationPolicyFunc(
		func(_ context.Context, _ string) error {
			calls++
			if status != "PENDING" {
				return domain.ErrVerificationNotPending
			}
			return nil
		},
	)

	key := uuid.NewString()
	verificationID := "ver_test_" + uuid.NewString()

	first, err := store.CreateDocumentIdempotent(
		ctx, key, "same_hash", verificationID, "SYNTHETIC_ID",
	)
	if err != nil {
		t.Fatal(err)
	}

	status = "APPROVED"

	replay, err := store.CreateDocumentIdempotent(
		ctx, key, "same_hash", verificationID, "SYNTHETIC_ID",
	)
	if err != nil {
		t.Fatal(err)
	}

	if !replay.Replayed || replay.ResourceID != first.ResourceID {
		t.Fatalf("incorrect replay: %+v", replay)
	}
	if calls != 1 {
		t.Fatalf("replay should not repeat verification lookup, got %d calls", calls)
	}

	_, err = store.CreateDocumentIdempotent(
		ctx, uuid.NewString(), "new_hash", verificationID, "SYNTHETIC_ID",
	)
	if !errors.Is(err, domain.ErrVerificationNotPending) {
		t.Fatalf("expected approval policy conflict, got %v", err)
	}

	_, err = store.CreateDocumentIdempotent(
		ctx, key, "different_hash", verificationID, "SYNTHETIC_ID",
	)
	if !errors.Is(err, idempotency.ErrConflict) {
		t.Fatalf("expected idempotency conflict, got %v", err)
	}
}
