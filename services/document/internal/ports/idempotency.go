package ports

import (
	"context"

	"idenova/idempotency"
)

type IdempotencyRepository interface {
	CreateDocumentIdempotent(
		ctx context.Context,
		key string,
		requestHash string,
		verificationID string,
		documentType string,
	) (idempotency.Result, error)
}
