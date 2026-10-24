package http

import (
	"encoding/json"
	"net/http"

	"idenova/identity/internal/application"
)

type Handler struct {
	create *application.CreateOrganization
	get    *application.GetOrganization
}

func NewHandler(
	create *application.CreateOrganization,
	get *application.GetOrganization,
) *Handler {
	return &Handler{
		create: create,
		get:    get,
	}
}

func (h *Handler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	org, err := h.create.Execute(r.Context(), request.ID, request.Name)
	if err != nil {
		handleDomainError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(org)
}

func (h *Handler) GetOrganization(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	org, err := h.get.Execute(r.Context(), id)
	if err != nil {
		handleDomainError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(org)
}
