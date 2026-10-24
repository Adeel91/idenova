package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"idenova/document/internal/domain"
)

func handleDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidDocument):
		writeError(w, http.StatusBadRequest, "Invalid document")

	case errors.Is(err, domain.ErrDocumentNotFound):
		writeError(w, http.StatusNotFound, "Document not found")

	case errors.Is(err, domain.ErrDocumentDuplicate):
		writeError(w, http.StatusConflict, "Document already exists")

	case errors.Is(err, domain.ErrInvalidTransition):
		writeError(w, http.StatusConflict, "Invalid document status transition")

	default:
		log.Printf("unexpected document error: %v", err)
		writeError(w, http.StatusInternalServerError, "Internal server error")
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("failed to encode JSON response: %v", err)
	}
}
