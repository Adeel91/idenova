package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"idenova/verification/internal/domain"
)

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	message := "Internal server error"

	switch {
	case errors.Is(err, domain.ErrNotFound):
		status = http.StatusNotFound
		message = "Verification not found"

	case errors.Is(err, domain.ErrDuplicate):
		status = http.StatusConflict
		message = "Verification already exists"

	case errors.Is(err, domain.ErrInvalidInput),
		errors.Is(err, domain.ErrInvalidStatus):
		status = http.StatusBadRequest
		message = "Invalid request"

	case errors.Is(err, domain.ErrInvalidTransition):
		status = http.StatusConflict
		message = "Invalid verification status transition"

	default:
		log.Printf("unexpected verification error: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if encodeErr := json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	}); encodeErr != nil {
		log.Printf("failed to encode error response: %v", encodeErr)
	}
}
