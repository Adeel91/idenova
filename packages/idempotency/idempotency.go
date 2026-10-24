package idempotency

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

var (
	ErrConflict   = errors.New("idempotency key reused with different request")
	ErrInvalidKey = errors.New("invalid idempotency key")
)

type Result struct {
	ResourceID string
	Replayed   bool
}

func ValidateKey(key string) error {
	if key == "" || len(key) > 255 || strings.TrimSpace(key) != key {
		return ErrInvalidKey
	}

	for _, character := range key {
		if character < 33 || character > 126 {
			return ErrInvalidKey
		}
	}

	return nil
}

func Fingerprint(operation, payload string) string {
	hash := sha256.Sum256([]byte(operation + ":" + payload))
	return hex.EncodeToString(hash[:])
}
