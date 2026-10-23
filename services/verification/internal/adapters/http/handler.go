package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"idenova/verification/internal/application"
	"idenova/verification/internal/domain"
)

type Handler struct {
	create *application.CreateVerification
	get    *application.GetVerification
	update *application.UpdateVerification
	list   *application.ListVerifications
}

func NewHandler(
	create *application.CreateVerification,
	get *application.GetVerification,
	list *application.ListVerifications,
	update *application.UpdateVerification,
) *Handler {
	return &Handler{create: create, get: get, update: update, list: list}
}

func (h *Handler) CreateVerification(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ID         string `json:"id"`
		CustomerID string `json:"customerId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, domain.ErrInvalidInput)
		return
	}

	err := h.create.Execute(r.Context(), request.ID, request.CustomerID)
	if err != nil {
		writeError(w, err)
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
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(verification); err != nil {
		log.Printf("failed to encode verification response: %v", err)
	}
}

func (h *Handler) UpdateVerification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var request struct {
		Status domain.VerificationStatus `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, domain.ErrInvalidInput)
		return
	}

	if request.Status != domain.StatusApproved &&
		request.Status != domain.StatusRejected {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	if err := h.update.Execute(r.Context(), id, request.Status); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListVerifications(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	limit := int32(10)
	offset := int32(0)

	if value := query.Get("limit"); value != "" {
		n, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			writeError(w, domain.ErrInvalidInput)
			return
		}
		limit = int32(n)
	}

	if value := query.Get("offset"); value != "" {
		n, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			writeError(w, domain.ErrInvalidInput)
			return
		}
		offset = int32(n)
	}

	var status *domain.VerificationStatus

	if value := query.Get("status"); value != "" {
		s := domain.VerificationStatus(value)
		status = &s
	}

	results, err := h.list.Execute(r.Context(), status, limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Printf("failed to encode verifications: %v", err)
	}
}
