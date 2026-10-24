package application

import (
	"context"
	"errors"
	"testing"

	"idenova/identity/internal/domain"
)

func TestGetOrganization(t *testing.T) {
	t.Run("existing organization", func(t *testing.T) {
		expected := &domain.Organization{
			ID:   "org_123",
			Name: "Acme GmbH",
		}

		repo := &fakeOrganizationRepository{found: expected}
		uc := NewGetOrganization(repo)

		org, err := uc.Execute(context.Background(), "org_123")
		if err != nil {
			t.Fatal(err)
		}

		if org == nil || org.ID != expected.ID || org.Name != expected.Name {
			t.Fatalf("unexpected organization: %+v", org)
		}
	})

	t.Run("organization not found", func(t *testing.T) {
		repo := &fakeOrganizationRepository{}
		uc := NewGetOrganization(repo)

		org, err := uc.Execute(context.Background(), "missing")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if org != nil {
			t.Fatalf("expected nil organization, got %+v", org)
		}
	})

	t.Run("repository failure", func(t *testing.T) {
		dbErr := errors.New("database unavailable")
		repo := &fakeOrganizationRepository{findErr: dbErr}
		uc := NewGetOrganization(repo)

		org, err := uc.Execute(context.Background(), "org_123")

		if !errors.Is(err, dbErr) {
			t.Fatalf("expected repository error, got %v", err)
		}

		if org != nil {
			t.Fatal("expected nil organization")
		}
	})
}
