package postgres

import (
	"context"
	"errors"

	"idenova/verification/internal/adapters/postgres/db"
	"idenova/verification/internal/domain"

	"github.com/jackc/pgx/v5/pgtype"
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

func (r *Repository) Update(ctx context.Context, v *domain.Verification) error {
	rows, err := r.queries.UpdateVerificationStatus(
		ctx,
		db.UpdateVerificationStatusParams{
			ID:     v.ID,
			Status: string(v.Status),
		},
	)
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("verification not found or no longer pending")
	}

	return nil
}

func (r *Repository) List(
	ctx context.Context,
	status *domain.VerificationStatus,
	limit, offset int32,
) ([]domain.Verification, error) {
	params := db.ListVerificationsParams{
		Limit:  limit,
		Offset: offset,
	}

	if status != nil {
		params.Status = pgtype.Text{
			String: string(*status),
			Valid:  true,
		}
	}

	rows, err := r.queries.ListVerifications(ctx, params)
	if err != nil {
		return nil, err
	}

	verifications := make([]domain.Verification, 0, len(rows))

	for _, row := range rows {
		verifications = append(verifications, domain.Verification{
			ID:         row.ID,
			CustomerID: row.CustomerID,
			Status:     domain.VerificationStatus(row.Status),
		})
	}

	return verifications, nil
}
