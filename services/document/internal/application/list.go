package application

import (
	"context"
	"strings"

	"idenova/document/internal/domain"
	"idenova/document/internal/ports"
)

type ListDocuments struct {
	repo ports.DocumentRepository
}

func NewListDocuments(repo ports.DocumentRepository) *ListDocuments {
	return &ListDocuments{repo: repo}
}

func (uc *ListDocuments) Execute(
	ctx context.Context,
	verificationID string,
) ([]domain.Document, error) {
	if strings.TrimSpace(verificationID) == "" {
		return nil, domain.ErrInvalidDocument
	}

	return uc.repo.ListByVerification(ctx, verificationID)
}
