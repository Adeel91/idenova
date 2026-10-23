package domain

import "errors"

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
		return errors.New("only pending verifications can be updated")
	}

	if status != StatusApproved && status != StatusRejected {
		return errors.New("invalid verification status")
	}

	v.Status = status
	return nil
}
