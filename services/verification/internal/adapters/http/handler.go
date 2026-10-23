package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Adeel91/idenova/services/verification/internal/application"
)

type Handler struct {
	create *application.CreateVerification
	get    *application.GetVerification
}

func NewHandler(
	create *application.CreateVerification,
	get *application.GetVerification,
) *Handler {
	return &Handler{create: create, get: get}
}

func (h *Handler) CreateVerification(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ID         string `json:"id"`
		CustomerID string `json:"customerId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.create.Execute(r.Context(), request.ID, request.CustomerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"id":     request.ID,
		"status": "PENDING",
	}); err != nil {
		log.Printf("failed to encode create response: %v", err)
	}
}

func (h *Handler) GetVerification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	verification, err := h.get.Execute(r.Context(), id)
	if err != nil {
		http.Error(w, "Verification not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(verification); err != nil {
		log.Printf("failed to encode verification response: %v", err)
	}
}
