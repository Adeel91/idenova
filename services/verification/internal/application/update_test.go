package application

import (
	"context"
	"testing"

	"idenova/verification/internal/adapters/memory"
	"idenova/verification/internal/domain"
)

func TestUpdateVerification(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewRepository()

	verification := &domain.Verification{
		ID:         "123",
		CustomerID: "456",
		Status:     domain.StatusPending,
	}

	if err := repo.Save(ctx, verification); err != nil {
		t.Fatal(err)
	}

	useCase := NewUpdateVerification(repo)

	if err := useCase.Execute(ctx, "123", domain.StatusApproved); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	saved, err := repo.FindByID(ctx, "123")
	if err != nil {
		t.Fatal(err)
	}

	if saved.Status != domain.StatusApproved {
		t.Errorf("expected APPROVED, got %s", saved.Status)
	}
}

func TestUpdateVerificationRejectsInvalidStatus(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewRepository()

	verification := &domain.Verification{
		ID:         "123",
		CustomerID: "456",
		Status:     domain.StatusPending,
	}

	if err := repo.Save(ctx, verification); err != nil {
		t.Fatal(err)
	}

	useCase := NewUpdateVerification(repo)

	err := useCase.Execute(ctx, "123", domain.VerificationStatus("INVALID"))
	if err == nil {
		t.Fatal("expected an error for invalid status")
	}

	saved, err := repo.FindByID(ctx, "123")
	if err != nil {
		t.Fatal(err)
	}

	if saved.Status != domain.StatusPending {
		t.Errorf("expected PENDING, got %s", saved.Status)
	}
}
