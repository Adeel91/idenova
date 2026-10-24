package domain

import "strings"

type UserRole string

const (
	RoleAdmin  UserRole = "ADMIN"
	RoleMember UserRole = "MEMBER"
)

type User struct {
	ID             string
	OrganizationID string
	Email          string
	Role           UserRole
}

func NewUser(
	id, organizationID, email string,
	role UserRole,
) (*User, error) {
	if strings.TrimSpace(id) == "" ||
		strings.TrimSpace(organizationID) == "" ||
		strings.TrimSpace(email) == "" {
		return nil, ErrInvalidUser
	}

	if !strings.Contains(email, "@") {
		return nil, ErrInvalidUser
	}

	if role != RoleAdmin && role != RoleMember {
		return nil, ErrInvalidUser
	}

	return &User{
		ID:             id,
		OrganizationID: organizationID,
		Email:          email,
		Role:           role,
	}, nil
}
