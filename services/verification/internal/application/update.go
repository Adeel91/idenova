package application

import (
	"context"

	"idenova/verification/internal/domain"
	"idenova/verification/internal/ports"
)

type UpdateVerification struct {
	repo ports.VerificationRepository
}

func NewUpdateVerification(repo ports.VerificationRepository) *UpdateVerification {
	return &UpdateVerification{repo: repo}
}

func (uc *UpdateVerification) Execute(
	ctx context.Context,
	id string,
	status domain.VerificationStatus,
) error {
	verification, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := verification.UpdateStatus(status); err != nil {
		return err
	}

	return uc.repo.Update(ctx, verification)
}
