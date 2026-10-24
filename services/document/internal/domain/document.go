package domain

import "strings"

type DocumentStatus string

const (
	StatusPending  DocumentStatus = "PENDING"
	StatusApproved DocumentStatus = "APPROVED"
	StatusRejected DocumentStatus = "REJECTED"
)

type Document struct {
	ID             string
	VerificationID string
	Type           string
	Status         DocumentStatus
}

func NewDocument(id, verificationID, documentType string) (*Document, error) {
	if strings.TrimSpace(id) == "" ||
		strings.TrimSpace(verificationID) == "" ||
		strings.TrimSpace(documentType) == "" {
		return nil, ErrInvalidDocument
	}

	return &Document{
		ID:             id,
		VerificationID: verificationID,
		Type:           documentType,
		Status:         StatusPending,
	}, nil
}

func (d *Document) UpdateStatus(status DocumentStatus) error {
	if d.Status != StatusPending {
		return ErrInvalidTransition
	}

	if status != StatusApproved && status != StatusRejected {
		return ErrInvalidDocument
	}

	d.Status = status
	return nil
}
