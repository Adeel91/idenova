package domain

import "errors"

var (
	ErrNotFound          = errors.New("verification not found")
	ErrDuplicate         = errors.New("verification already exists")
	ErrInvalidInput      = errors.New("invalid verification input")
	ErrInvalidStatus     = errors.New("invalid verification status")
	ErrInvalidTransition = errors.New("invalid verification status transition")
)
