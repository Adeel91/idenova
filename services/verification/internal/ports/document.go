package ports

import "context"

type DocumentStatus string

const (
	DocumentPending  DocumentStatus = "PENDING"
	DocumentApproved DocumentStatus = "APPROVED"
	DocumentRejected DocumentStatus = "REJECTED"
)

type Document struct {
	ID     string
	Status DocumentStatus
}

type DocumentReader interface {
	ListByVerificationID(
		ctx context.Context,
		verificationID string,
	) ([]Document, error)
}
