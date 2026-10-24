package http

import (
	"encoding/json"
	"net/http"

	"idenova/identity/internal/application"
	"idenova/identity/internal/domain"
)

type UserHandler struct {
	create *application.CreateUser
	get    *application.GetUser
	list   *application.ListUsers
}

func NewUserHandler(
	create *application.CreateUser,
	get *application.GetUser,
	list *application.ListUsers,
) *UserHandler {
	return &UserHandler{
		create: create,
		get:    get,
		list:   list,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ID             string          `json:"id"`
		OrganizationID string          `json:"organizationId"`
		Email          string          `json:"email"`
		Role           domain.UserRole `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	user, err := h.create.Execute(
		r.Context(),
		request.ID,
		request.OrganizationID,
		request.Email,
		request.Role,
	)
	if err != nil {
		handleDomainError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		return
	}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.get.Execute(r.Context(), r.PathValue("id"))
	if err != nil {
		handleDomainError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.list.Execute(r.Context(), r.PathValue("id"))
	if err != nil {
		handleDomainError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(users)
}
