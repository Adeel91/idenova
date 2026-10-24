package idempotency

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateKey(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{"valid UUID", "c6bb4738-7a34-4f96-8393-d8beb63a4561", false},
		{"valid custom key", "request_123", false},
		{"empty key", "", true},
		{"whitespace only", "   ", true},
		{"leading whitespace", " key123", true},
		{"trailing whitespace", "key123 ", true},
		{"newline", "key\n123", true},
		{"maximum length", strings.Repeat("a", 255), false},
		{"too long", strings.Repeat("a", 256), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateKey(tt.key)

			if tt.wantErr && !errors.Is(err, ErrInvalidKey) {
				t.Fatalf("expected ErrInvalidKey, got %v", err)
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestFingerprint(t *testing.T) {
	first := Fingerprint("create_verification", "customer_123")
	repeated := Fingerprint("create_verification", "customer_123")
	differentCustomer := Fingerprint("create_verification", "customer_456")
	differentOperation := Fingerprint("create_document", "customer_123")

	if first != repeated {
		t.Fatal("same input must produce the same fingerprint")
	}

	if first == differentCustomer {
		t.Fatal("different requests must produce different fingerprints")
	}

	if first == differentOperation {
		t.Fatal("different operations must produce different fingerprints")
	}

	if len(first) != 64 {
		t.Fatalf("expected 64-character SHA-256 hash, got %d", len(first))
	}
}
