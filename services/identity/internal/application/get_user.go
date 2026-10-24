package application

import (
	"context"

	"idenova/identity/internal/domain"
	"idenova/identity/internal/ports"
)

type GetUser struct {
	repo ports.UserRepository
}

func NewGetUser(repo ports.UserRepository) *GetUser {
	return &GetUser{repo: repo}
}

func (uc *GetUser) Execute(
	ctx context.Context,
	id string,
) (*domain.User, error) {
	return uc.repo.FindByID(ctx, id)
}
