package ports

import (
	"context"

	"github.com/Adeel91/idenova/services/verification/internal/domain"
)

type VerificationRepository interface {
	FindByID(ctx context.Context, id string) (*domain.Verification, error)
	Save(ctx context.Context, verification *domain.Verification) error
}
