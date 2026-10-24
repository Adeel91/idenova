package domain

import "errors"

var (
	ErrInvalidOrganization   = errors.New("invalid organization")
	ErrOrganizationNotFound  = errors.New("organization not found")
	ErrOrganizationDuplicate = errors.New("organization already exists")
)
