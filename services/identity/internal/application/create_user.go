package application

import (
	"context"

	"idenova/identity/internal/domain"
	"idenova/identity/internal/ports"
)

type CreateUser struct {
	repo ports.UserRepository
}

func NewCreateUser(repo ports.UserRepository) *CreateUser {
	return &CreateUser{repo: repo}
}

func (uc *CreateUser) Execute(
	ctx context.Context,
	id, organizationID, email string,
	role domain.UserRole,
) (*domain.User, error) {
	user, err := domain.NewUser(id, organizationID, email, role)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
