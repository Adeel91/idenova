package domain

import (
	"strings"
)

type Organization struct {
	ID   string
	Name string
}

func NewOrganization(id, name string) (*Organization, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(name) == "" {
		return nil, ErrInvalidOrganization
	}

	return &Organization{
		ID:   id,
		Name: name,
	}, nil
}
