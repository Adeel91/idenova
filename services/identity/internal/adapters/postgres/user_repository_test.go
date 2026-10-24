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

func TestUserRepositoryIntegration(t *testing.T) {
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

	queries := db.New(pool)
	orgRepo := NewRepository(queries)
	userRepo := NewUserRepository(queries)

	org := &domain.Organization{
		ID:   "org_" + uuid.NewString(),
		Name: "Test Organization",
	}

	if err := orgRepo.Save(ctx, org); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_, err := pool.Exec(
			context.Background(),
			"DELETE FROM organizations WHERE id = $1",
			org.ID,
		)
		if err != nil {
			t.Errorf("cleanup failed: %v", err)
		}
	})

	user := &domain.User{
		ID:             "user_" + uuid.NewString(),
		OrganizationID: org.ID,
		Email:          "alice@example.com",
		Role:           domain.RoleAdmin,
	}

	t.Run("create user", func(t *testing.T) {
		if err := userRepo.Save(ctx, user); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("retrieve user", func(t *testing.T) {
		found, err := userRepo.FindByID(ctx, user.ID)
		if err != nil {
			t.Fatal(err)
		}

		if found.ID != user.ID ||
			found.OrganizationID != org.ID ||
			found.Email != user.Email ||
			found.Role != user.Role {
			t.Fatalf("unexpected user: %+v", found)
		}
	})

	t.Run("duplicate user ID", func(t *testing.T) {
		err := userRepo.Save(ctx, user)

		if !errors.Is(err, domain.ErrUserDuplicate) {
			t.Fatalf("expected duplicate error, got %v", err)
		}
	})

	t.Run("duplicate organization email", func(t *testing.T) {
		duplicate := &domain.User{
			ID:             "user_" + uuid.NewString(),
			OrganizationID: org.ID,
			Email:          user.Email,
			Role:           domain.RoleMember,
		}

		err := userRepo.Save(ctx, duplicate)

		if !errors.Is(err, domain.ErrUserDuplicate) {
			t.Fatalf("expected duplicate error, got %v", err)
		}
	})

	t.Run("missing organization", func(t *testing.T) {
		missing := &domain.User{
			ID:             "user_" + uuid.NewString(),
			OrganizationID: "missing_" + uuid.NewString(),
			Email:          "missing@example.com",
			Role:           domain.RoleMember,
		}

		err := userRepo.Save(ctx, missing)

		if !errors.Is(err, domain.ErrOrganizationNotFound) {
			t.Fatalf("expected organization not found, got %v", err)
		}
	})

	t.Run("list organization users", func(t *testing.T) {
		users, err := userRepo.ListByOrganization(ctx, org.ID)
		if err != nil {
			t.Fatal(err)
		}

		if len(users) != 1 || users[0].ID != user.ID {
			t.Fatalf("unexpected users: %+v", users)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		_, err := userRepo.FindByID(ctx, "missing_"+uuid.NewString())

		if !errors.Is(err, domain.ErrUserNotFound) {
			t.Fatalf("expected user not found, got %v", err)
		}
	})
}
