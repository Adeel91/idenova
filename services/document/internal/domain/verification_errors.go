package domain

import "errors"

var (
	ErrVerificationNotFound   = errors.New("verification not found")
	ErrVerificationNotPending = errors.New("verification is not pending")
)
