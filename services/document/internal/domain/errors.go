package domain

import "errors"

var (
	ErrInvalidDocument   = errors.New("invalid document")
	ErrDocumentNotFound  = errors.New("document not found")
	ErrDocumentDuplicate = errors.New("document already exists")
	ErrInvalidTransition = errors.New("invalid document status transition")
)
