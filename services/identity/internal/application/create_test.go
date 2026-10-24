package application

import (
	"context"
	"errors"
	"testing"

	"idenova/identity/internal/domain"
)

type fakeOrganizationRepository struct {
	saved   *domain.Organization
	saveErr error
	found   *domain.Organization
	findErr error
}

func (r *fakeOrganizationRepository) Save(
	_ context.Context,
	org *domain.Organization,
) error {
	if r.saveErr != nil {
		return r.saveErr
	}

	r.saved = org
	return nil
}

func (r *fakeOrganizationRepository) FindByID(
	_ context.Context,
	_ string,
) (*domain.Organization, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}

	return r.found, nil
}

func TestCreateOrganization(t *testing.T) {
	t.Run("valid organization", func(t *testing.T) {
		repo := &fakeOrganizationRepository{}
		uc := NewCreateOrganization(repo)

		org, err := uc.Execute(
			context.Background(),
			"org_123",
			"Acme GmbH",
		)
		if err != nil {
			t.Fatal(err)
		}

		if org.ID != "org_123" || org.Name != "Acme GmbH" {
			t.Fatalf("unexpected organization: %+v", org)
		}

		if repo.saved == nil || repo.saved.ID != org.ID {
			t.Fatal("organization was not saved")
		}
	})

	t.Run("invalid organization", func(t *testing.T) {
		repo := &fakeOrganizationRepository{}
		uc := NewCreateOrganization(repo)

		org, err := uc.Execute(context.Background(), "", "Acme GmbH")

		if !errors.Is(err, domain.ErrInvalidOrganization) {
			t.Fatalf("expected validation error, got %v", err)
		}

		if org != nil || repo.saved != nil {
			t.Fatal("invalid organization should not be saved")
		}
	})

	t.Run("repository failure", func(t *testing.T) {
		dbErr := errors.New("database unavailable")
		repo := &fakeOrganizationRepository{saveErr: dbErr}
		uc := NewCreateOrganization(repo)

		org, err := uc.Execute(
			context.Background(),
			"org_123",
			"Acme GmbH",
		)

		if !errors.Is(err, dbErr) {
			t.Fatalf("expected repository error, got %v", err)
		}

		if org != nil {
			t.Fatal("expected nil organization on failure")
		}
	})
}
