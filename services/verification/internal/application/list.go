package application

import (
	"context"

	"idenova/verification/internal/domain"
	"idenova/verification/internal/ports"
)

type ListVerifications struct {
	repo ports.VerificationRepository
}

func NewListVerifications(repo ports.VerificationRepository) *ListVerifications {
	return &ListVerifications{repo: repo}
}

func (uc *ListVerifications) Execute(
	ctx context.Context,
	status *domain.VerificationStatus,
	limit, offset int32,
) ([]domain.Verification, error) {
	if limit < 1 || limit > 100 || offset < 0 {
		return nil, domain.ErrInvalidInput
	}

	if status != nil &&
		*status != domain.StatusPending &&
		*status != domain.StatusApproved &&
		*status != domain.StatusRejected {
		return nil, domain.ErrInvalidStatus
	}

	return uc.repo.List(ctx, status, limit, offset)
}
