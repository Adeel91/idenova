package application

import (
	"context"

	"idenova/identity/internal/domain"
	"idenova/identity/internal/ports"
)

type ListUsers struct {
	repo ports.UserRepository
}

func NewListUsers(repo ports.UserRepository) *ListUsers {
	return &ListUsers{repo: repo}
}

func (uc *ListUsers) Execute(
	ctx context.Context,
	organizationID string,
) ([]domain.User, error) {
	if organizationID == "" {
		return nil, domain.ErrInvalidUser
	}

	return uc.repo.ListByOrganization(ctx, organizationID)
}
