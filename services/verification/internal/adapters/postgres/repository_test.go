package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"idenova/verification/internal/adapters/postgres/db"
	"idenova/verification/internal/domain"

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
