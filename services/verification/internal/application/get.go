package application

import (
	"context"

	"idenova/verification/internal/domain"
	"idenova/verification/internal/ports"
)

type GetVerification struct {
	repo ports.VerificationRepository
}

func NewGetVerification(repo ports.VerificationRepository) *GetVerification {
	return &GetVerification{repo: repo}
}

func (uc *GetVerification) Execute(ctx context.Context, id string) (*domain.Verification, error) {
	return uc.repo.FindByID(ctx, id)
}
