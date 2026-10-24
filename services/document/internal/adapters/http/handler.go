package http

import (
	"encoding/json"
	"net/http"

	"idenova/document/internal/application"
	"idenova/document/internal/domain"
)

type Handler struct {
	createIdempotent *application.CreateDocumentIdempotent
	create           *application.CreateDocument
	get              *application.GetDocument
	list             *application.ListDocuments
	update           *application.UpdateDocument
}

func NewHandler(
	create *application.CreateDocument,
	get *application.GetDocument,
	list *application.ListDocuments,
	update *application.UpdateDocument,
) *Handler {
	return &Handler{
		create: create,
		get:    get,
		list:   list,
		update: update,
	}
}

func (h *Handler) SetIdempotentCreation(
	useCase *application.CreateDocumentIdempotent,
) {
	h.createIdempotent = useCase
}

func (h *Handler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ID             string `json:"id"`
		VerificationID string `json:"verificationId"`
		Type           string `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if h.createIdempotent != nil &&
		r.Header.Get("Idempotency-Key") == "" {
		writeError(w, http.StatusBadRequest, "Idempotency-Key header is required")
		return
	}

	if key := r.Header.Get("Idempotency-Key"); key != "" {
		if h.createIdempotent == nil {
			writeError(w, http.StatusInternalServerError, "Idempotent creation unavailable")
			return
		}

		result, err := h.createIdempotent.Execute(
			r.Context(),
			key,
			request.VerificationID,
			request.Type,
		)
		if err != nil {
			handleDomainError(w, err)
			return
		}

		document, err := h.get.Execute(r.Context(), result.ResourceID)
		if err != nil {
			handleDomainError(w, err)
			return
		}

		if result.Replayed {
			w.Header().Set("Idempotency-Replayed", "true")
			writeJSON(w, http.StatusOK, document)
			return
		}

		writeJSON(w, http.StatusCreated, document)
		return
	}

	document, err := h.create.Execute(
		r.Context(),
		request.ID,
		request.VerificationID,
		request.Type,
	)
	if err != nil {
		handleDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, document)
}

func (h *Handler) GetDocument(w http.ResponseWriter, r *http.Request) {
	document, err := h.get.Execute(r.Context(), r.PathValue("id"))
	if err != nil {
		handleDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, document)
}

func (h *Handler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	documents, err := h.list.Execute(r.Context(), r.PathValue("id"))
	if err != nil {
		handleDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, documents)
}

func (h *Handler) UpdateDocument(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Status domain.DocumentStatus `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if err := h.update.Execute(
		r.Context(),
		r.PathValue("id"),
		request.Status,
	); err != nil {
		handleDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
