package postgres

import (
	"context"

	"idenova/verification/internal/adapters/postgres/db"
	"idenova/verification/internal/domain"
)

type Repository struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) *Repository {
	return &Repository{queries: queries}
}

func (r *Repository) Save(ctx context.Context, v *domain.Verification) error {
	return r.queries.CreateVerification(ctx, db.CreateVerificationParams{
		ID:         v.ID,
		CustomerID: v.CustomerID,
		Status:     string(v.Status),
	})
}

func (r *Repository) FindByID(ctx context.Context, id string) (*domain.Verification, error) {
	result, err := r.queries.GetVerification(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.Verification{
		ID:         result.ID,
		CustomerID: result.CustomerID,
		Status:     domain.VerificationStatus(result.Status),
	}, nil
}
