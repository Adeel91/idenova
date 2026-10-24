package application

import (
	"context"
	"errors"
	"testing"

	"idenova/identity/internal/domain"
)

type fakeUserRepository struct {
	users   map[string]*domain.User
	saveErr error
	findErr error
	listErr error
}

func (r *fakeUserRepository) Save(
	_ context.Context,
	user *domain.User,
) error {
	if r.saveErr != nil {
		return r.saveErr
	}

	if _, exists := r.users[user.ID]; exists {
		return domain.ErrUserDuplicate
	}

	copy := *user
	r.users[user.ID] = &copy
	return nil
}

func (r *fakeUserRepository) FindByID(
	_ context.Context,
	id string,
) (*domain.User, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}

	user, exists := r.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	copy := *user
	return &copy, nil
}

func (r *fakeUserRepository) ListByOrganization(
	_ context.Context,
	organizationID string,
) ([]domain.User, error) {
	if r.listErr != nil {
		return nil, r.listErr
	}

	results := make([]domain.User, 0)

	for _, user := range r.users {
		if user.OrganizationID == organizationID {
			results = append(results, *user)
		}
	}

	return results, nil
}

func TestCreateUser(t *testing.T) {
	t.Run("valid user", func(t *testing.T) {
		repo := &fakeUserRepository{
			users: make(map[string]*domain.User),
		}
		uc := NewCreateUser(repo)

		user, err := uc.Execute(
			context.Background(),
			"user_123",
			"org_123",
			"alice@example.com",
			domain.RoleAdmin,
		)
		if err != nil {
			t.Fatal(err)
		}

		if user.ID != "user_123" || user.Role != domain.RoleAdmin {
			t.Fatalf("unexpected user: %+v", user)
		}

		if _, exists := repo.users[user.ID]; !exists {
			t.Fatal("user was not saved")
		}
	})

	t.Run("invalid user", func(t *testing.T) {
		repo := &fakeUserRepository{
			users: make(map[string]*domain.User),
		}
		uc := NewCreateUser(repo)

		_, err := uc.Execute(
			context.Background(),
			"user_123",
			"org_123",
			"invalid-email",
			domain.RoleMember,
		)

		if !errors.Is(err, domain.ErrInvalidUser) {
			t.Fatalf("expected invalid user error, got %v", err)
		}

		if len(repo.users) != 0 {
			t.Fatal("invalid user should not be saved")
		}
	})

	t.Run("duplicate user", func(t *testing.T) {
		repo := &fakeUserRepository{
			users: make(map[string]*domain.User),
		}
		uc := NewCreateUser(repo)

		for i := 0; i < 2; i++ {
			_, err := uc.Execute(
				context.Background(),
				"user_123",
				"org_123",
				"alice@example.com",
				domain.RoleAdmin,
			)

			if i == 1 && !errors.Is(err, domain.ErrUserDuplicate) {
				t.Fatalf("expected duplicate error, got %v", err)
			}
		}
	})

	t.Run("repository failure", func(t *testing.T) {
		dbErr := errors.New("database unavailable")
		repo := &fakeUserRepository{
			saveErr: dbErr,
		}
		uc := NewCreateUser(repo)

		user, err := uc.Execute(
			context.Background(),
			"user_123",
			"org_123",
			"alice@example.com",
			domain.RoleMember,
		)

		if !errors.Is(err, dbErr) || user != nil {
			t.Fatalf("expected repository error, got user=%v err=%v", user, err)
		}
	})
}

func TestGetUser(t *testing.T) {
	repo := &fakeUserRepository{
		users: map[string]*domain.User{
			"user_123": {
				ID:             "user_123",
				OrganizationID: "org_123",
				Email:          "alice@example.com",
				Role:           domain.RoleAdmin,
			},
		},
	}

	uc := NewGetUser(repo)

	t.Run("existing user", func(t *testing.T) {
		user, err := uc.Execute(context.Background(), "user_123")
		if err != nil {
			t.Fatal(err)
		}

		if user.ID != "user_123" {
			t.Fatalf("unexpected user: %+v", user)
		}
	})

	t.Run("missing user", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), "missing")

		if !errors.Is(err, domain.ErrUserNotFound) {
			t.Fatalf("expected not found error, got %v", err)
		}
	})
}

func TestListUsers(t *testing.T) {
	repo := &fakeUserRepository{
		users: map[string]*domain.User{
			"user_1": {
				ID:             "user_1",
				OrganizationID: "org_123",
				Email:          "alice@example.com",
				Role:           domain.RoleAdmin,
			},
			"user_2": {
				ID:             "user_2",
				OrganizationID: "org_123",
				Email:          "bob@example.com",
				Role:           domain.RoleMember,
			},
			"user_3": {
				ID:             "user_3",
				OrganizationID: "org_456",
				Email:          "charlie@example.com",
				Role:           domain.RoleMember,
			},
		},
	}

	uc := NewListUsers(repo)

	t.Run("organization users", func(t *testing.T) {
		users, err := uc.Execute(context.Background(), "org_123")
		if err != nil {
			t.Fatal(err)
		}

		if len(users) != 2 {
			t.Fatalf("expected 2 users, got %d", len(users))
		}
	})

	t.Run("invalid organization ID", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), "")

		if !errors.Is(err, domain.ErrInvalidUser) {
			t.Fatalf("expected invalid user error, got %v", err)
		}
	})
}
