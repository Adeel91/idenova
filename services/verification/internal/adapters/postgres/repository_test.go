package postgres

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"idenova/verification/internal/adapters/postgres/db"
	"idenova/verification/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestRepositoryIntegration(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		pool.Close()
	})

	if err := pool.Ping(ctx); err != nil {
		t.Fatal(err)
	}

	repo := NewRepository(db.New(pool))

	id := "integration_" + time.Now().Format("20060102150405.000000000")

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(
			context.Background(), 5*time.Second,
		)
		defer cleanupCancel()

		if _, err := pool.Exec(
			cleanupCtx,
			"DELETE FROM verifications WHERE id = $1",
			id,
		); err != nil {
			t.Errorf("cleanup failed: %v", err)
		}
	})

	// CREATE
	verification := &domain.Verification{
		ID:         id,
		CustomerID: "customer123",
		Status:     domain.StatusPending,
	}

	if err := repo.Save(ctx, verification); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	err = repo.Save(ctx, verification)

	if !errors.Is(err, domain.ErrDuplicate) {
		t.Fatalf("expected ErrDuplicate, got %v", err)
	}

	// GET
	saved, err := repo.FindByID(ctx, id)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if saved.CustomerID != "customer123" ||
		saved.Status != domain.StatusPending {
		t.Fatalf("unexpected saved verification: %+v", saved)
	}

	// UPDATE
	if err := saved.UpdateStatus(domain.StatusApproved); err != nil {
		t.Fatal(err)
	}

	if err := repo.Update(ctx, saved); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	updated, err := repo.FindByID(ctx, id)
	if err != nil {
		t.Fatal(err)
	}

	if updated.Status != domain.StatusApproved {
		t.Fatalf("expected APPROVED, got %s", updated.Status)
	}

	// REJECT SECOND UPDATE
	updated.Status = domain.StatusRejected

	if err := repo.Update(ctx, updated); err == nil {
		t.Fatal("expected second update to fail")
	}

	// LIST WITH FILTER
	status := domain.StatusApproved

	results, err := repo.List(ctx, &status, 100, 0)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}

	found := false
	for _, item := range results {
		if item.ID == id {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("verification missing from filtered list")
	}

	// GET MISSING RECORD
	_, err = repo.FindByID(ctx, id+"_missing")
	if err == nil {
		t.Fatal("expected error for missing verification")
	}
}

func TestConcurrentStatusUpdates(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	queries := db.New(pool)
	repo := NewRepository(queries)

	verification := &domain.Verification{
		ID:         "concurrent-test-" + uuid.NewString(),
		CustomerID: "customer-123",
		Status:     domain.StatusPending,
	}

	if err := repo.Save(ctx, verification); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(),
			"DELETE FROM verifications WHERE id = $1",
			verification.ID,
		)
	})

	var wg sync.WaitGroup
	results := make(chan error, 2)

	for _, status := range []domain.VerificationStatus{
		domain.StatusApproved,
		domain.StatusRejected,
	} {
		wg.Add(1)

		go func(status domain.VerificationStatus) {
			defer wg.Done()

			updated := *verification
			updated.Status = status

			results <- repo.Update(ctx, &updated)
		}(status)
	}

	wg.Wait()
	close(results)

	successCount := 0
	conflictCount := 0

	for err := range results {
		switch {
		case err == nil:
			successCount++
		case errors.Is(err, domain.ErrInvalidTransition):
			conflictCount++
		default:
			t.Fatalf("unexpected update error: %v", err)
		}
	}

	if successCount != 1 || conflictCount != 1 {
		t.Fatalf(
			"expected one success and one conflict, got %d successes and %d conflicts",
			successCount,
			conflictCount,
		)
	}
}
