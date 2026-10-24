package application

import (
	"context"
	"strings"

	"idenova/idempotency"
	"idenova/verification/internal/domain"
	"idenova/verification/internal/ports"
)

var ErrInvalidIdempotencyKey = idempotency.ErrInvalidKey

type CreateVerificationIdempotent struct {
	repo ports.IdempotencyRepository
}

func NewCreateVerificationIdempotent(
	repo ports.IdempotencyRepository,
) *CreateVerificationIdempotent {
	return &CreateVerificationIdempotent{repo: repo}
}

func (uc *CreateVerificationIdempotent) Execute(
	ctx context.Context,
	key string,
	customerID string,
) (ports.IdempotencyResult, error) {
	if err := idempotency.ValidateKey(key); err != nil {
		return ports.IdempotencyResult{}, err
	}

	customerID = strings.TrimSpace(customerID)

	if customerID == "" || len(customerID) > 255 {
		return ports.IdempotencyResult{}, domain.ErrInvalidInput
	}

	requestHash := idempotency.Fingerprint(
		"create_verification",
		customerID,
	)

	return uc.repo.CreateVerificationIdempotent(
		ctx,
		key,
		requestHash,
		customerID,
	)
}
