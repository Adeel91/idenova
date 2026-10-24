package postgres

import (
	"context"
	"errors"
	"os"
	"testing"

	"idenova/identity/internal/adapters/postgres/db"
	"idenova/identity/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestOrganizationRepositoryIntegration(t *testing.T) {
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

	org := &domain.Organization{
		ID:   "org_" + uuid.NewString(),
		Name: "Acme GmbH",
	}

	t.Cleanup(func() {
		_, err := pool.Exec(
			context.Background(),
			"DELETE FROM organizations WHERE id = $1",
			org.ID,
		)
		if err != nil {
			t.Errorf("failed to clean up test organization: %v", err)
		}
	})

	t.Run("create organization", func(t *testing.T) {
		if err := repo.Save(ctx, org); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("retrieve organization", func(t *testing.T) {
		found, err := repo.FindByID(ctx, org.ID)
		if err != nil {
			t.Fatal(err)
		}

		if found.ID != org.ID || found.Name != org.Name {
			t.Fatalf("unexpected organization: %+v", found)
		}
	})

	t.Run("duplicate organization", func(t *testing.T) {
		err := repo.Save(ctx, org)

		if !errors.Is(err, domain.ErrOrganizationDuplicate) {
			t.Fatalf("expected duplicate error, got %v", err)
		}
	})

	t.Run("organization not found", func(t *testing.T) {
		_, err := repo.FindByID(ctx, "missing_"+uuid.NewString())

		if !errors.Is(err, domain.ErrOrganizationNotFound) {
			t.Fatalf("expected not found error, got %v", err)
		}
	})
}
