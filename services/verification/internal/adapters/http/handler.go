package http

import (
	"encoding/json"
	"net/http"

	"github.com/Adeel91/idenova/services/verification/internal/application"
)

type Handler struct {
	create *application.CreateVerification
}

func NewHandler(create *application.CreateVerification) *Handler {
	return &Handler{create: create}
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

	json.NewEncoder(w).Encode(map[string]string{
		"id":     request.ID,
		"status": "PENDING",
	})
}
