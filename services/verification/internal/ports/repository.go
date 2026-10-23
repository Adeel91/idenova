package ports

import (
	"context"

	"idenova/verification/internal/domain"
)

type VerificationRepository interface {
	FindByID(ctx context.Context, id string) (*domain.Verification, error)
	Save(ctx context.Context, verification *domain.Verification) error
	Update(ctx context.Context, verification *domain.Verification) error

	List(
		ctx context.Context,
		status *domain.VerificationStatus,
		limit, offset int32,
	) ([]domain.Verification, error)
}
