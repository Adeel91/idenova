package ports

import (
	"context"

	"idenova/idempotency"
)

var ErrIdempotencyConflict = idempotency.ErrConflict

type IdempotencyResult = idempotency.Result

type IdempotencyRepository interface {
	CreateVerificationIdempotent(
		ctx context.Context,
		key string,
		requestHash string,
		customerID string,
	) (IdempotencyResult, error)
}
