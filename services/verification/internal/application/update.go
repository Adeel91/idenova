package application

import (
	"context"
	"errors"

	"idenova/verification/internal/domain"
	"idenova/verification/internal/ports"
)

var ErrDocumentPolicy = errors.New("document requirements not satisfied")

type UpdateVerification struct {
	repo      ports.VerificationRepository
	documents ports.DocumentReader
}

func NewUpdateVerification(
	repo ports.VerificationRepository,
	documents ports.DocumentReader,
) *UpdateVerification {
	return &UpdateVerification{
		repo:      repo,
		documents: documents,
	}
}

func (uc *UpdateVerification) Execute(
	ctx context.Context,
	id string,
	status domain.VerificationStatus,
) error {
	verification, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if verification.Status != domain.StatusPending {
		return domain.ErrInvalidTransition
	}

	if status != domain.StatusApproved && status != domain.StatusRejected {
		return domain.ErrInvalidStatus
	}

	if uc.documents == nil {
		return errors.New("document reader not configured")
	}

	documents, err := uc.documents.ListByVerificationID(ctx, id)
	if err != nil {
		return err
	}

	switch status {
	case domain.StatusApproved:
		if len(documents) == 0 {
			return ErrDocumentPolicy
		}

		for _, document := range documents {
			if document.Status != ports.DocumentApproved {
				return ErrDocumentPolicy
			}
		}

	case domain.StatusRejected:
		hasRejectedDocument := false

		for _, document := range documents {
			if document.Status == ports.DocumentRejected {
				hasRejectedDocument = true
				break
			}
		}

		if !hasRejectedDocument {
			return ErrDocumentPolicy
		}
	}

	if err := verification.UpdateStatus(status); err != nil {
		return err
	}

	return uc.repo.Update(ctx, verification)
}
