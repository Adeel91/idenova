package domain

import (
	"errors"
	"testing"
)

func TestNewOrganization(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		orgName string
		wantErr bool
	}{
		{
			name:    "valid organization",
			id:      "org_123",
			orgName: "Acme GmbH",
		},
		{
			name:    "empty ID",
			id:      "",
			orgName: "Acme GmbH",
			wantErr: true,
		},
		{
			name:    "empty name",
			id:      "org_123",
			orgName: "",
			wantErr: true,
		},
		{
			name:    "whitespace ID",
			id:      "   ",
			orgName: "Acme GmbH",
			wantErr: true,
		},
		{
			name:    "whitespace name",
			id:      "org_123",
			orgName: "   ",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			org, err := NewOrganization(tt.id, tt.orgName)

			if tt.wantErr {
				if !errors.Is(err, ErrInvalidOrganization) {
					t.Fatalf("expected ErrInvalidOrganization, got %v", err)
				}

				if org != nil {
					t.Fatal("expected nil organization")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if org == nil {
				t.Fatal("expected organization, got nil")
			}

			if org.ID != tt.id || org.Name != tt.orgName {
				t.Errorf("unexpected organization: %+v", org)
			}
		})
	}
}
