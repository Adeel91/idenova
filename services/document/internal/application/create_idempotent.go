package application

import (
	"context"
	"encoding/json"
	"strings"

	"idenova/document/internal/domain"
	"idenova/document/internal/ports"
	"idenova/idempotency"
)

type CreateDocumentIdempotent struct {
	repo ports.IdempotencyRepository
}

func NewCreateDocumentIdempotent(
	repo ports.IdempotencyRepository,
) *CreateDocumentIdempotent {
	return &CreateDocumentIdempotent{repo: repo}
}

func (uc *CreateDocumentIdempotent) Execute(
	ctx context.Context,
	key string,
	verificationID string,
	documentType string,
) (idempotency.Result, error) {
	if err := idempotency.ValidateKey(key); err != nil {
		return idempotency.Result{}, err
	}

	verificationID = strings.TrimSpace(verificationID)
	documentType = strings.TrimSpace(documentType)

	if verificationID == "" || documentType == "" ||
		len(verificationID) > 255 || len(documentType) > 255 {
		return idempotency.Result{}, domain.ErrInvalidDocument
	}

	payload, err := json.Marshal(struct {
		VerificationID string `json:"verificationId"`
		Type           string `json:"type"`
	}{
		VerificationID: verificationID,
		Type:           documentType,
	})
	if err != nil {
		return idempotency.Result{}, err
	}

	hash := idempotency.Fingerprint(
		"create_document",
		string(payload),
	)

	return uc.repo.CreateDocumentIdempotent(
		ctx,
		key,
		hash,
		verificationID,
		documentType,
	)
}
