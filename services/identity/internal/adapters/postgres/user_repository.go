package postgres

import (
	"context"
	"errors"

	"idenova/identity/internal/adapters/postgres/db"
	"idenova/identity/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepository struct {
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) *UserRepository {
	return &UserRepository{queries: queries}
}

func (r *UserRepository) Save(
	ctx context.Context,
	user *domain.User,
) error {
	err := r.queries.CreateUser(ctx, db.CreateUserParams{
		ID:             user.ID,
		OrganizationID: user.OrganizationID,
		Email:          user.Email,
		Role:           string(user.Role),
	})

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return domain.ErrUserDuplicate
			case "23503":
				return domain.ErrOrganizationNotFound
			}
		}

		return err
	}

	return nil
}

func (r *UserRepository) FindByID(
	ctx context.Context,
	id string,
) (*domain.User, error) {
	record, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	return &domain.User{
		ID:             record.ID,
		OrganizationID: record.OrganizationID,
		Email:          record.Email,
		Role:           domain.UserRole(record.Role),
	}, nil
}

func (r *UserRepository) ListByOrganization(
	ctx context.Context,
	organizationID string,
) ([]domain.User, error) {
	records, err := r.queries.ListUsersByOrganization(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	users := make([]domain.User, 0, len(records))

	for _, record := range records {
		users = append(users, domain.User{
			ID:             record.ID,
			OrganizationID: record.OrganizationID,
			Email:          record.Email,
			Role:           domain.UserRole(record.Role),
		})
	}

	return users, nil
}
