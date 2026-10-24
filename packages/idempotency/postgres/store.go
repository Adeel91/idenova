package postgres

import (
	"context"
	"crypto/rand"
	"fmt"

	"idenova/idempotency"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

type CreateFunc func(
	ctx context.Context,
	tx pgx.Tx,
	resourceID string,
) error

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) Execute(
	ctx context.Context,
	operation string,
	key string,
	requestHash string,
	idPrefix string,
	create CreateFunc,
) (idempotency.Result, error) {
	if err := idempotency.ValidateKey(key); err != nil {
		return idempotency.Result{}, err
	}

	if operation == "" || requestHash == "" || idPrefix == "" || create == nil {
		return idempotency.Result{}, fmt.Errorf("invalid idempotency configuration")
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return idempotency.Result{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	resourceID, err := generateResourceID(idPrefix)
	if err != nil {
		return idempotency.Result{}, err
	}

	tag, err := tx.Exec(
		ctx,
		`INSERT INTO idempotency_keys (
			operation, idempotency_key, request_hash, resource_id
		) VALUES ($1, $2, $3, $4)
		ON CONFLICT (operation, idempotency_key) DO NOTHING`,
		operation,
		key,
		requestHash,
		resourceID,
	)
	if err != nil {
		return idempotency.Result{}, err
	}

	if tag.RowsAffected() == 0 {
		var existingHash string
		var existingID string

		err = tx.QueryRow(
			ctx,
			`SELECT request_hash, resource_id
			 FROM idempotency_keys
			 WHERE operation = $1 AND idempotency_key = $2`,
			operation,
			key,
		).Scan(&existingHash, &existingID)

		if err != nil {
			return idempotency.Result{}, err
		}

		if existingHash != requestHash {
			return idempotency.Result{}, idempotency.ErrConflict
		}

		if err := tx.Commit(ctx); err != nil {
			return idempotency.Result{}, err
		}

		return idempotency.Result{
			ResourceID: existingID,
			Replayed:   true,
		}, nil
	}

	if err := create(ctx, tx, resourceID); err != nil {
		return idempotency.Result{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return idempotency.Result{}, err
	}

	return idempotency.Result{
		ResourceID: resourceID,
		Replayed:   false,
	}, nil
}

func generateResourceID(prefix string) (string, error) {
	var id [16]byte

	if _, err := rand.Read(id[:]); err != nil {
		return "", err
	}

	id[6] = (id[6] & 0x0f) | 0x40
	id[8] = (id[8] & 0x3f) | 0x80

	return fmt.Sprintf(
		"%s%x-%x-%x-%x-%x",
		prefix,
		id[0:4],
		id[4:6],
		id[6:8],
		id[8:10],
		id[10:16],
	), nil
}
