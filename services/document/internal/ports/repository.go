package ports

import (
	"context"

	"idenova/document/internal/domain"
)

type DocumentRepository interface {
	Save(ctx context.Context, document *domain.Document) error
	FindByID(ctx context.Context, id string) (*domain.Document, error)
	ListByVerification(
		ctx context.Context,
		verificationID string,
	) ([]domain.Document, error)
	Update(ctx context.Context, document *domain.Document) error
}
