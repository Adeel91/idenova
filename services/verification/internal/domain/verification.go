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

func (v *Verification) UpdateStatus(status VerificationStatus) error {
	if v.Status != StatusPending {
		return ErrInvalidTransition
	}

	if status != StatusApproved && status != StatusRejected {
		return ErrInvalidStatus
	}

	v.Status = status
	return nil
}
