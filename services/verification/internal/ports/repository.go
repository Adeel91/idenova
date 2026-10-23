package ports

import (
	"context"

	"idenova/verification/internal/domain"
)

type VerificationRepository interface {
	FindByID(ctx context.Context, id string) (*domain.Verification, error)
	Save(ctx context.Context, verification *domain.Verification) error
	Update(ctx context.Context, verification *domain.Verification) error
}
