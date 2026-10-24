package application

import (
	"context"

	"idenova/identity/internal/domain"
	"idenova/identity/internal/ports"
)

type CreateOrganization struct {
	repo ports.OrganizationRepository
}

func NewCreateOrganization(repo ports.OrganizationRepository) *CreateOrganization {
	return &CreateOrganization{repo: repo}
}

func (uc *CreateOrganization) Execute(
	ctx context.Context,
	id, name string,
) (*domain.Organization, error) {
	org, err := domain.NewOrganization(id, name)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, org); err != nil {
		return nil, err
	}

	return org, nil
}
