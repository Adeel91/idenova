package ports

import (
	"context"

	"github.com/Adeel91/idenova/services/verification/internal/domain"
)

type VerificationRepository interface {
	Save(ctx context.Context, verification *domain.Verification) error
}
