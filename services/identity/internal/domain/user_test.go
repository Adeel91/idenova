package domain

import (
	"errors"
	"testing"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		orgID   string
		email   string
		role    UserRole
		wantErr bool
	}{
		{
			name:  "valid admin",
			id:    "user_123",
			orgID: "org_123",
			email: "alice@example.com",
			role:  RoleAdmin,
		},
		{
			name:  "valid member",
			id:    "user_456",
			orgID: "org_123",
			email: "bob@example.com",
			role:  RoleMember,
		},
		{
			name:    "empty user ID",
			id:      "",
			orgID:   "org_123",
			email:   "alice@example.com",
			role:    RoleAdmin,
			wantErr: true,
		},
		{
			name:    "empty organization ID",
			id:      "user_123",
			orgID:   "",
			email:   "alice@example.com",
			role:    RoleMember,
			wantErr: true,
		},
		{
			name:    "invalid email",
			id:      "user_123",
			orgID:   "org_123",
			email:   "invalid-email",
			role:    RoleMember,
			wantErr: true,
		},
		{
			name:    "empty email",
			id:      "user_123",
			orgID:   "org_123",
			email:   "",
			role:    RoleMember,
			wantErr: true,
		},
		{
			name:    "invalid role",
			id:      "user_123",
			orgID:   "org_123",
			email:   "alice@example.com",
			role:    UserRole("OWNER"),
			wantErr: true,
		},
		{
			name:    "whitespace user ID",
			id:      "   ",
			orgID:   "org_123",
			email:   "alice@example.com",
			role:    RoleMember,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.id, tt.orgID, tt.email, tt.role)

			if tt.wantErr {
				if !errors.Is(err, ErrInvalidUser) {
					t.Fatalf("expected ErrInvalidUser, got %v", err)
				}

				if user != nil {
					t.Fatalf("expected nil user, got %+v", user)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if user == nil {
				t.Fatal("expected user, got nil")
			}

			if user.ID != tt.id ||
				user.OrganizationID != tt.orgID ||
				user.Email != tt.email ||
				user.Role != tt.role {
				t.Fatalf("unexpected user: %+v", user)
			}
		})
	}
}
