package postgres

import (
	"context"

	"idenova/document/internal/domain"
	"idenova/idempotency"
	idempotencypg "idenova/idempotency/postgres"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const createDocumentOperation = "create_document"

type VerificationReader interface {
	RequirePending(ctx context.Context, verificationID string) error
}

type IdempotencyStore struct {
	store    *idempotencypg.Store
	verifier VerificationReader
}

func NewIdempotencyStore(
	pool *pgxpool.Pool,
	verifier VerificationReader,
) *IdempotencyStore {
	if verifier == nil {
		panic("verification reader is required")
	}

	return &IdempotencyStore{
		store:    idempotencypg.NewStore(pool),
		verifier: verifier,
	}
}

func (s *IdempotencyStore) CreateDocumentIdempotent(
	ctx context.Context,
	key string,
	requestHash string,
	verificationID string,
	documentType string,
) (idempotency.Result, error) {
	return s.store.Execute(
		ctx,
		createDocumentOperation,
		key,
		requestHash,
		"doc_",
		func(ctx context.Context, tx pgx.Tx, resourceID string) error {
			document, err := domain.NewDocument(
				resourceID,
				verificationID,
				documentType,
			)
			if err != nil {
				return err
			}

			if err := s.verifier.RequirePending(
				ctx,
				document.VerificationID,
			); err != nil {
				return err
			}

			_, err = tx.Exec(
				ctx,
				`INSERT INTO documents (
					id, verification_id, document_type, status
				) VALUES ($1, $2, $3, $4)`,
				document.ID,
				document.VerificationID,
				document.Type,
				string(document.Status),
			)
			return err
		},
	)
}
