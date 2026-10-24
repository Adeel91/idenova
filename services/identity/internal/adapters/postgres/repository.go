package postgres

import (
	"context"
	"errors"

	"idenova/identity/internal/adapters/postgres/db"
	"idenova/identity/internal/domain"

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
	org *domain.Organization,
) error {
	err := r.queries.CreateOrganization(ctx, db.CreateOrganizationParams{
		ID:   org.ID,
		Name: org.Name,
	})
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrOrganizationDuplicate
		}

		return err
	}

	return nil
}

func (r *Repository) FindByID(
	ctx context.Context,
	id string,
) (*domain.Organization, error) {
	record, err := r.queries.GetOrganizationByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOrganizationNotFound
		}
		return nil, err
	}

	return &domain.Organization{
		ID:   record.ID,
		Name: record.Name,
	}, nil
}
