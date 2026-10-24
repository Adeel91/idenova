package application

import (
	"context"

	"idenova/document/internal/domain"
	"idenova/document/internal/ports"
)

type UpdateDocument struct {
	repo ports.DocumentRepository
}

func NewUpdateDocument(repo ports.DocumentRepository) *UpdateDocument {
	return &UpdateDocument{repo: repo}
}

func (uc *UpdateDocument) Execute(
	ctx context.Context,
	id string,
	status domain.DocumentStatus,
) error {
	document, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := document.UpdateStatus(status); err != nil {
		return err
	}

	return uc.repo.Update(ctx, document)
}
