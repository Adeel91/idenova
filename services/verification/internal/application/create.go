package application

import (
	"context"
	"errors"

	"github.com/Adeel91/idenova/services/verification/internal/domain"
	"github.com/Adeel91/idenova/services/verification/internal/ports"
)

type CreateVerification struct {
	repo ports.VerificationRepository
}

func NewCreateVerification(repo ports.VerificationRepository) *CreateVerification {
	return &CreateVerification{repo: repo}
}

func (uc *CreateVerification) Execute(ctx context.Context, id, customerID string) error {
	if id == "" || customerID == "" {
		return errors.New("verification ID and customer ID are required")
	}

	verification := &domain.Verification{
		ID:         id,
		CustomerID: customerID,
		Status:     domain.StatusPending,
	}

	return uc.repo.Save(ctx, verification)
}
