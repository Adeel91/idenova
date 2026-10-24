package application

import (
	"context"

	"idenova/document/internal/domain"
	"idenova/document/internal/ports"
)

type GetDocument struct {
	repo ports.DocumentRepository
}

func NewGetDocument(repo ports.DocumentRepository) *GetDocument {
	return &GetDocument{repo: repo}
}

func (uc *GetDocument) Execute(ctx context.Context, id string) (*domain.Document, error) {
	return uc.repo.FindByID(ctx, id)
}
