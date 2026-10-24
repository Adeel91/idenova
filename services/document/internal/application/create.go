package application

import (
	"context"

	"idenova/document/internal/domain"
	"idenova/document/internal/ports"
)

type CreateDocument struct {
	repo ports.DocumentRepository
}

func NewCreateDocument(repo ports.DocumentRepository) *CreateDocument {
	return &CreateDocument{repo: repo}
}

func (uc *CreateDocument) Execute(
	ctx context.Context,
	id, verificationID, documentType string,
) (*domain.Document, error) {
	document, err := domain.NewDocument(id, verificationID, documentType)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, document); err != nil {
		return nil, err
	}

	return document, nil
}
