package domain

import "errors"

var (
	ErrInvalidOrganization   = errors.New("invalid organization")
	ErrOrganizationNotFound  = errors.New("organization not found")
	ErrOrganizationDuplicate = errors.New("organization already exists")
	ErrInvalidUser           = errors.New("invalid user")
	ErrUserNotFound          = errors.New("user not found")
	ErrUserDuplicate         = errors.New("user already exists")
)
