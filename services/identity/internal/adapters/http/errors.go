package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"idenova/identity/internal/domain"
)

func handleDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidOrganization):
		writeError(w, http.StatusBadRequest, "Invalid organization")

	case errors.Is(err, domain.ErrOrganizationNotFound):
		writeError(w, http.StatusNotFound, "Organization not found")

	case errors.Is(err, domain.ErrOrganizationDuplicate):
		writeError(w, http.StatusConflict, "Organization already exists")

	case errors.Is(err, domain.ErrInvalidUser):
		writeError(w, http.StatusBadRequest, "Invalid user")

	case errors.Is(err, domain.ErrUserNotFound):
		writeError(w, http.StatusNotFound, "User not found")

	case errors.Is(err, domain.ErrUserDuplicate):
		writeError(w, http.StatusConflict, "User already exists")

	default:
		log.Printf("unexpected identity error: %v", err)
		writeError(w, http.StatusInternalServerError, "Internal server error")
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	}); err != nil {
		log.Printf("failed to encode error response: %v", err)
	}
}
