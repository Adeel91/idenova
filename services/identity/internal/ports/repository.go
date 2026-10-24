package ports

import (
	"context"

	"idenova/identity/internal/domain"
)

type OrganizationRepository interface {
	Save(ctx context.Context, org *domain.Organization) error
	FindByID(ctx context.Context, id string) (*domain.Organization, error)
}
