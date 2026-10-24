package ports

import (
	"context"

	"idenova/identity/internal/domain"
)

type OrganizationRepository interface {
	Save(ctx context.Context, org *domain.Organization) error
	FindByID(ctx context.Context, id string) (*domain.Organization, error)
}

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
	ListByOrganization(
		ctx context.Context,
		organizationID string,
	) ([]domain.User, error)
}
