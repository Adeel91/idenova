package postgres

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"

	"idenova/document/internal/adapters/postgres/db"
	"idenova/document/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestDocumentRepositoryIntegration(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)

	if err := pool.Ping(ctx); err != nil {
		t.Fatal(err)
	}

	repo := NewRepository(db.New(pool))

	document := &domain.Document{
		ID:             "doc_" + uuid.NewString(),
		VerificationID: "verification_" + uuid.NewString(),
		Type:           "SYNTHETIC_ID",
		Status:         domain.StatusPending,
	}

	t.Cleanup(func() {
		_, err := pool.Exec(
			context.Background(),
			"DELETE FROM documents WHERE id = $1",
			document.ID,
		)
		if err != nil {
			t.Errorf("cleanup failed: %v", err)
		}
	})

	t.Run("create document", func(t *testing.T) {
		if err := repo.Save(ctx, document); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("retrieve document", func(t *testing.T) {
		found, err := repo.FindByID(ctx, document.ID)
		if err != nil {
			t.Fatal(err)
		}

		if found.ID != document.ID ||
			found.VerificationID != document.VerificationID ||
			found.Type != document.Type ||
			found.Status != domain.StatusPending {
			t.Fatalf("unexpected document: %+v", found)
		}
	})

	t.Run("duplicate document", func(t *testing.T) {
		err := repo.Save(ctx, document)

		if !errors.Is(err, domain.ErrDocumentDuplicate) {
			t.Fatalf("expected duplicate error, got %v", err)
		}
	})

	t.Run("list by verification", func(t *testing.T) {
		documents, err := repo.ListByVerification(ctx, document.VerificationID)
		if err != nil {
			t.Fatal(err)
		}

		if len(documents) != 1 || documents[0].ID != document.ID {
			t.Fatalf("unexpected documents: %+v", documents)
		}
	})

	t.Run("missing document", func(t *testing.T) {
		_, err := repo.FindByID(ctx, "missing_"+uuid.NewString())

		if !errors.Is(err, domain.ErrDocumentNotFound) {
			t.Fatalf("expected not found error, got %v", err)
		}
	})

	t.Run("concurrent status updates", func(t *testing.T) {
		var wg sync.WaitGroup
		results := make(chan error, 2)

		for _, status := range []domain.DocumentStatus{
			domain.StatusApproved,
			domain.StatusRejected,
		} {
			wg.Add(1)

			go func(status domain.DocumentStatus) {
				defer wg.Done()

				updated := *document
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
				t.Fatalf("unexpected error: %v", err)
			}
		}

		if successCount != 1 || conflictCount != 1 {
			t.Fatalf(
				"expected one success and one conflict, got %d successes and %d conflicts",
				successCount,
				conflictCount,
			)
		}

		found, err := repo.FindByID(ctx, document.ID)
		if err != nil {
			t.Fatal(err)
		}

		if found.Status != domain.StatusApproved &&
			found.Status != domain.StatusRejected {
			t.Fatalf("unexpected final status: %s", found.Status)
		}
	})
}
