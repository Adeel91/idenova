package domain

type VerificationStatus string

const (
	StatusPending  VerificationStatus = "PENDING"
	StatusApproved VerificationStatus = "APPROVED"
	StatusRejected VerificationStatus = "REJECTED"
)

type Verification struct {
	ID         string
	CustomerID string
	Status     VerificationStatus
}
