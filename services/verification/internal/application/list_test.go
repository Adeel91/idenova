package application

import (
	"context"
	"testing"

	"idenova/verification/internal/adapters/memory"
	"idenova/verification/internal/domain"
)

func TestListVerifications(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewRepository()

	records := []domain.Verification{
		{ID: "001", CustomerID: "alice", Status: domain.StatusPending},
		{ID: "002", CustomerID: "bob", Status: domain.StatusApproved},
		{ID: "003", CustomerID: "charlie", Status: domain.StatusPending},
	}

	for i := range records {
		if err := repo.Save(ctx, &records[i]); err != nil {
			t.Fatal(err)
		}
	}

	useCase := NewListVerifications(repo)
	status := domain.StatusPending

	results, err := useCase.Execute(ctx, &status, 1, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 verification, got %d", len(results))
	}

	if results[0].ID != "003" {
		t.Errorf("expected verification 003, got %s", results[0].ID)
	}
}

func TestListVerificationsInvalidInput(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewRepository()
	useCase := NewListVerifications(repo)

	tests := []struct {
		name   string
		status *domain.VerificationStatus
		limit  int32
		offset int32
	}{
		{
			name:  "zero limit",
			limit: 0,
		},
		{
			name:   "negative offset",
			limit:  10,
			offset: -1,
		},
		{
			name:  "limit too large",
			limit: 101,
		},
		{
			name:   "invalid status",
			status: statusPtr(domain.VerificationStatus("UNKNOWN")),
			limit:  10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := useCase.Execute(ctx, tt.status, tt.limit, tt.offset)

			if err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func statusPtr(status domain.VerificationStatus) *domain.VerificationStatus {
	return &status
}
