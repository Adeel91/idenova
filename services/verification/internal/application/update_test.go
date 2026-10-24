package application

import (
	"context"
	"errors"
	"testing"

	"idenova/verification/internal/adapters/memory"
	"idenova/verification/internal/domain"
	"idenova/verification/internal/ports"
)

type fakeDocumentReader struct {
	documents []ports.Document
	err       error
}

func (f fakeDocumentReader) ListByVerificationID(
	_ context.Context,
	_ string,
) ([]ports.Document, error) {
	return f.documents, f.err
}

func TestUpdateVerification(t *testing.T) {
	tests := []struct {
		name       string
		initial    domain.VerificationStatus
		requested  domain.VerificationStatus
		documents  []ports.Document
		readerErr  error
		wantStatus domain.VerificationStatus
		wantErr    error
	}{
		{
			name:      "approve all approved",
			initial:   domain.StatusPending,
			requested: domain.StatusApproved,
			documents: []ports.Document{
				{ID: "doc1", Status: ports.DocumentApproved},
				{ID: "doc2", Status: ports.DocumentApproved},
			},
			wantStatus: domain.StatusApproved,
		},
		{
			name:       "cannot approve without documents",
			initial:    domain.StatusPending,
			requested:  domain.StatusApproved,
			wantStatus: domain.StatusPending,
			wantErr:    ErrDocumentPolicy,
		},
		{
			name:      "cannot approve pending document",
			initial:   domain.StatusPending,
			requested: domain.StatusApproved,
			documents: []ports.Document{
				{ID: "doc1", Status: ports.DocumentApproved},
				{ID: "doc2", Status: ports.DocumentPending},
			},
			wantStatus: domain.StatusPending,
			wantErr:    ErrDocumentPolicy,
		},
		{
			name:      "cannot approve rejected document",
			initial:   domain.StatusPending,
			requested: domain.StatusApproved,
			documents: []ports.Document{
				{ID: "doc1", Status: ports.DocumentRejected},
			},
			wantStatus: domain.StatusPending,
			wantErr:    ErrDocumentPolicy,
		},
		{
			name:      "reject with rejected document",
			initial:   domain.StatusPending,
			requested: domain.StatusRejected,
			documents: []ports.Document{
				{ID: "doc1", Status: ports.DocumentApproved},
				{ID: "doc2", Status: ports.DocumentRejected},
			},
			wantStatus: domain.StatusRejected,
		},
		{
			name:       "cannot reject without rejected document",
			initial:    domain.StatusPending,
			requested:  domain.StatusRejected,
			documents:  []ports.Document{{ID: "doc1", Status: ports.DocumentApproved}},
			wantStatus: domain.StatusPending,
			wantErr:    ErrDocumentPolicy,
		},
		{
			name:       "reject already approved verification",
			initial:    domain.StatusApproved,
			requested:  domain.StatusRejected,
			wantStatus: domain.StatusApproved,
			wantErr:    domain.ErrInvalidTransition,
		},
		{
			name:       "invalid status",
			initial:    domain.StatusPending,
			requested:  domain.VerificationStatus("UNKNOWN"),
			wantStatus: domain.StatusPending,
			wantErr:    domain.ErrInvalidStatus,
		},
		{
			name:       "document service error",
			initial:    domain.StatusPending,
			requested:  domain.StatusApproved,
			readerErr:  errors.New("document service unavailable"),
			wantStatus: domain.StatusPending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			repo := memory.NewRepository()

			err := repo.Save(ctx, &domain.Verification{
				ID:         "123",
				CustomerID: "456",
				Status:     tt.initial,
			})
			if err != nil {
				t.Fatal(err)
			}

			reader := fakeDocumentReader{
				documents: tt.documents,
				err:       tt.readerErr,
			}

			useCase := NewUpdateVerification(repo, reader)

			err = useCase.Execute(ctx, "123", tt.requested)

			switch {
			case tt.wantErr != nil && !errors.Is(err, tt.wantErr):
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			case tt.readerErr != nil && !errors.Is(err, tt.readerErr):
				t.Fatalf("expected reader error %v, got %v", tt.readerErr, err)
			case tt.wantErr == nil && tt.readerErr == nil && err != nil:
				t.Fatalf("unexpected error: %v", err)
			}

			saved, err := repo.FindByID(ctx, "123")
			if err != nil {
				t.Fatal(err)
			}

			if saved.Status != tt.wantStatus {
				t.Errorf(
					"expected %s, got %s",
					tt.wantStatus,
					saved.Status,
				)
			}
		})
	}
}
