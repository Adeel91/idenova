package postgres

import (
	"context"

	"idenova/idempotency"
	idempotencypg "idenova/idempotency/postgres"
	"idenova/verification/internal/ports"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrIdempotencyConflict = idempotency.ErrConflict

const createVerificationOperation = "create_verification"

type IdempotencyStore struct {
	store *idempotencypg.Store
}

func NewIdempotencyStore(pool *pgxpool.Pool) *IdempotencyStore {
	return &IdempotencyStore{
		store: idempotencypg.NewStore(pool),
	}
}

func (s *IdempotencyStore) CreateVerificationIdempotent(
	ctx context.Context,
	key string,
	requestHash string,
	customerID string,
) (ports.IdempotencyResult, error) {
	result, err := s.store.Execute(
		ctx,
		createVerificationOperation,
		key,
		requestHash,
		"ver_",
		func(ctx context.Context, tx pgx.Tx, resourceID string) error {
			_, err := tx.Exec(
				ctx,
				`INSERT INTO verifications (id, customer_id, status)
				 VALUES ($1, $2, 'PENDING')`,
				resourceID,
				customerID,
			)
			return err
		},
	)
	if err != nil {
		return ports.IdempotencyResult{}, err
	}

	return result, nil
}
