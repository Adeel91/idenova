package application

import (
	"context"

	"idenova/identity/internal/domain"
	"idenova/identity/internal/ports"
)

type GetOrganization struct {
	repo ports.OrganizationRepository
}

func NewGetOrganization(repo ports.OrganizationRepository) *GetOrganization {
	return &GetOrganization{repo: repo}
}

func (uc *GetOrganization) Execute(
	ctx context.Context,
	id string,
) (*domain.Organization, error) {
	return uc.repo.FindByID(ctx, id)
}
