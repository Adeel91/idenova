package postgres

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupIdempotencyTest(t *testing.T) (
	context.Context,
	*pgxpool.Pool,
	*IdempotencyStore,
	string,
) {
	t.Helper()

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	t.Cleanup(cancel)

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)

	if err := pool.Ping(ctx); err != nil {
		t.Fatal(err)
	}

	key := uuid.NewString()
	store := NewIdempotencyStore(pool)

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cleanupCancel()

		var resourceID string

		err := pool.QueryRow(
			cleanupCtx,
			`SELECT resource_id
			 FROM idempotency_keys
			 WHERE operation = $1 AND idempotency_key = $2`,
			createVerificationOperation,
			key,
		).Scan(&resourceID)

		if err == nil {
			if _, err := pool.Exec(
				cleanupCtx,
				"DELETE FROM idempotency_keys WHERE operation = $1 AND idempotency_key = $2",
				createVerificationOperation,
				key,
			); err != nil {
				t.Errorf("cleanup idempotency key: %v", err)
			}

			if _, err := pool.Exec(
				cleanupCtx,
				"DELETE FROM verifications WHERE id = $1",
				resourceID,
			); err != nil {
				t.Errorf("cleanup verification: %v", err)
			}
		}
	})

	return ctx, pool, store, key
}

func TestIdempotentCreationAndReplay(t *testing.T) {
	ctx, _, store, key := setupIdempotencyTest(t)

	first, err := store.CreateVerificationIdempotent(
		ctx, key, "hash_123", "customer_123",
	)
	if err != nil {
		t.Fatal(err)
	}

	if first.ResourceID == "" || first.Replayed {
		t.Fatalf("unexpected first result: %+v", first)
	}

	second, err := store.CreateVerificationIdempotent(
		ctx, key, "hash_123", "customer_123",
	)
	if err != nil {
		t.Fatal(err)
	}

	if second.ResourceID != first.ResourceID || !second.Replayed {
		t.Fatalf("unexpected replay result: %+v", second)
	}
}

func TestIdempotencyKeyConflict(t *testing.T) {
	ctx, _, store, key := setupIdempotencyTest(t)

	_, err := store.CreateVerificationIdempotent(
		ctx, key, "hash_original", "customer_123",
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.CreateVerificationIdempotent(
		ctx, key, "hash_different", "customer_456",
	)

	if !errors.Is(err, ErrIdempotencyConflict) {
		t.Fatalf("expected idempotency conflict, got %v", err)
	}
}

func TestConcurrentIdempotentCreation(t *testing.T) {
	ctx, pool, store, key := setupIdempotencyTest(t)

	const workers = 10

	type result struct {
		id       string
		replayed bool
		err      error
	}

	results := make(chan result, workers)
	start := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			<-start

			created, err := store.CreateVerificationIdempotent(
				ctx, key, "hash_concurrent", "customer_123",
			)

			results <- result{
				id:       created.ResourceID,
				replayed: created.Replayed,
				err:      err,
			}
		}()
	}

	close(start)
	wg.Wait()
	close(results)

	var resourceID string
	createdCount := 0
	replayedCount := 0

	for result := range results {
		if result.err != nil {
			t.Fatalf("concurrent request failed: %v", result.err)
		}

		if resourceID == "" {
			resourceID = result.id
		}

		if result.id != resourceID {
			t.Fatalf("different resource IDs: %s and %s", resourceID, result.id)
		}

		if result.replayed {
			replayedCount++
		} else {
			createdCount++
		}
	}

	if createdCount != 1 || replayedCount != workers-1 {
		t.Fatalf(
			"expected 1 creation and %d replays, got %d creations and %d replays",
			workers-1,
			createdCount,
			replayedCount,
		)
	}

	var count int

	err := pool.QueryRow(
		ctx,
		"SELECT COUNT(*) FROM verifications WHERE id = $1",
		resourceID,
	).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Fatalf("expected exactly one verification, got %d", count)
	}
}

func TestFailedTransactionDoesNotReserveKey(t *testing.T) {
	ctx, pool, store, key := setupIdempotencyTest(t)

	// customer_id is VARCHAR(255), so this insert must fail.
	invalidCustomerID := strings.Repeat("x", 256)

	_, err := store.CreateVerificationIdempotent(
		ctx, key, "hash_rollback", invalidCustomerID,
	)
	if err == nil {
		t.Fatal("expected database insertion error")
	}

	var count int

	err = pool.QueryRow(
		ctx,
		`SELECT COUNT(*) FROM idempotency_keys
		 WHERE operation = $1 AND idempotency_key = $2`,
		createVerificationOperation,
		key,
	).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}

	if count != 0 {
		t.Fatal("failed transaction left an idempotency record")
	}

	created, err := store.CreateVerificationIdempotent(
		ctx, key, "hash_rollback", "customer_valid",
	)
	if err != nil {
		t.Fatalf("retry after rollback failed: %v", err)
	}

	if created.Replayed || created.ResourceID == "" {
		t.Fatalf("unexpected retry result: %+v", created)
	}
}
