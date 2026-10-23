package domain

import "testing"

func TestUpdateStatus(t *testing.T) {
	v := &Verification{
		ID:         "123",
		CustomerID: "456",
		Status:     StatusPending,
	}

	err := v.UpdateStatus(StatusApproved)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v.Status != StatusApproved {
		t.Errorf("expected APPROVED, got %s", v.Status)
	}

	err = v.UpdateStatus(StatusRejected)

	if err == nil {
		t.Error("expected error when updating an approved verification")
	}
}

func TestUpdateStatusWithInvalidStatus(t *testing.T) {
	v := &Verification{
		ID:         "123",
		CustomerID: "456",
		Status:     StatusPending,
	}

	err := v.UpdateStatus(VerificationStatus("UNKNOWN"))

	if err == nil {
		t.Fatal("expected error for invalid status")
	}

	if v.Status != StatusPending {
		t.Errorf("expected PENDING, got %s", v.Status)
	}
}
