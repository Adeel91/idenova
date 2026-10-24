package postgres

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"

	"idenova/document/internal/domain"
	"idenova/idempotency"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func testIdempotencyStore(t *testing.T) *IdempotencyStore {
	t.Helper()

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(pool.Close)

	if err := pool.Ping(context.Background()); err != nil {
		t.Fatal(err)
	}

	return NewIdempotencyStore(pool, verificationPolicyFunc(
		func(context.Context, string) error { return nil },
	))
}

func TestIdempotentDocumentCreationAndReplay(t *testing.T) {
	store := testIdempotencyStore(t)
	ctx := context.Background()
	key := uuid.NewString()
	verificationID := "ver_test_" + uuid.NewString()

	first, err := store.CreateDocumentIdempotent(
		ctx, key, "hash_1", verificationID, "PASSPORT",
	)
	if err != nil {
		t.Fatal(err)
	}
	if first.Replayed || first.ResourceID == "" {
		t.Fatalf("unexpected first result: %+v", first)
	}

	second, err := store.CreateDocumentIdempotent(
		ctx, key, "hash_1", verificationID, "PASSPORT",
	)
	if err != nil {
		t.Fatal(err)
	}
	if !second.Replayed || second.ResourceID != first.ResourceID {
		t.Fatalf("unexpected replay result: %+v", second)
	}
}

func TestIdempotentDocumentConflict(t *testing.T) {
	store := testIdempotencyStore(t)
	ctx := context.Background()
	key := uuid.NewString()
	verificationID := "ver_test_" + uuid.NewString()

	_, err := store.CreateDocumentIdempotent(
		ctx, key, "original_hash", verificationID, "PASSPORT",
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.CreateDocumentIdempotent(
		ctx, key, "different_hash", verificationID, "PASSPORT",
	)
	if !errors.Is(err, idempotency.ErrConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
}

func TestConcurrentIdempotentDocumentCreation(t *testing.T) {
	store := testIdempotencyStore(t)
	ctx := context.Background()
	key := uuid.NewString()
	verificationID := "ver_test_" + uuid.NewString()

	const workers = 10

	results := make([]idempotency.Result, workers)
	errs := make([]error, workers)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			results[index], errs[index] = store.CreateDocumentIdempotent(
				ctx, key, "same_hash", verificationID, "PASSPORT",
			)
		}(i)
	}

	wg.Wait()

	created := 0
	replayed := 0

	for i := range results {
		if errs[i] != nil {
			t.Fatalf("worker %d: %v", i, errs[i])
		}
		if results[i].ResourceID != results[0].ResourceID {
			t.Fatalf("worker %d returned a different document ID", i)
		}
		if results[i].Replayed {
			replayed++
		} else {
			created++
		}
	}

	if created != 1 || replayed != workers-1 {
		t.Fatalf("expected 1 creation and 9 replays, got %d and %d", created, replayed)
	}
}

func TestFailedDocumentCreationDoesNotReserveKey(t *testing.T) {
	store := testIdempotencyStore(t)
	ctx := context.Background()
	key := uuid.NewString()

	_, err := store.CreateDocumentIdempotent(
		ctx, key, "invalid_hash", "", "PASSPORT",
	)
	if !errors.Is(err, domain.ErrInvalidDocument) {
		t.Fatalf("expected invalid document, got %v", err)
	}

	result, err := store.CreateDocumentIdempotent(
		ctx, key, "valid_hash", "ver_test_"+uuid.NewString(), "PASSPORT",
	)
	if err != nil {
		t.Fatal(err)
	}
	if result.Replayed {
		t.Fatal("failed creation must not reserve idempotency key")
	}
}
