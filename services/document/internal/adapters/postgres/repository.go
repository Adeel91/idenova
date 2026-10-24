package postgres

import (
	"context"
	"errors"

	"idenova/document/internal/adapters/postgres/db"
	"idenova/document/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) *Repository {
	return &Repository{queries: queries}
}

func (r *Repository) Save(
	ctx context.Context,
	document *domain.Document,
) error {
	err := r.queries.CreateDocument(ctx, db.CreateDocumentParams{
		ID:             document.ID,
		VerificationID: document.VerificationID,
		DocumentType:   document.Type,
		Status:         string(document.Status),
	})

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrDocumentDuplicate
		}

		return err
	}

	return nil
}

func (r *Repository) FindByID(
	ctx context.Context,
	id string,
) (*domain.Document, error) {
	record, err := r.queries.GetDocumentByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDocumentNotFound
		}

		return nil, err
	}

	return &domain.Document{
		ID:             record.ID,
		VerificationID: record.VerificationID,
		Type:           record.DocumentType,
		Status:         domain.DocumentStatus(record.Status),
	}, nil
}

func (r *Repository) ListByVerification(
	ctx context.Context,
	verificationID string,
) ([]domain.Document, error) {
	records, err := r.queries.ListDocumentsByVerification(ctx, verificationID)
	if err != nil {
		return nil, err
	}

	documents := make([]domain.Document, 0, len(records))

	for _, record := range records {
		documents = append(documents, domain.Document{
			ID:             record.ID,
			VerificationID: record.VerificationID,
			Type:           record.DocumentType,
			Status:         domain.DocumentStatus(record.Status),
		})
	}

	return documents, nil
}

func (r *Repository) Update(
	ctx context.Context,
	document *domain.Document,
) error {
	rows, err := r.queries.UpdateDocumentStatus(
		ctx,
		db.UpdateDocumentStatusParams{
			ID:     document.ID,
			Status: string(document.Status),
		},
	)
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrInvalidTransition
	}

	return nil
}
