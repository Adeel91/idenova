package application

import (
	"context"
	"testing"

	"github.com/Adeel91/idenova/services/verification/internal/adapters/memory"
	"github.com/Adeel91/idenova/services/verification/internal/domain"
)

func TestCreateVerification(t *testing.T) {
	repo := memory.NewRepository()
	useCase := NewCreateVerification(repo)

	err := useCase.Execute(
		context.Background(),
		"verification-123",
		"customer-456",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	saved, err := repo.FindByID(
		context.Background(),
		"verification-123",
	)

	if err != nil {
		t.Fatalf("verification was not saved: %v", err)
	}

	if saved.ID != "verification-123" {
		t.Errorf("unexpected verification ID: %s", saved.ID)
	}

	if saved.CustomerID != "customer-456" {
		t.Errorf("unexpected customer ID: %s", saved.CustomerID)
	}

	if saved.Status != domain.StatusPending {
		t.Errorf("expected PENDING, got %s", saved.Status)
	}

}
