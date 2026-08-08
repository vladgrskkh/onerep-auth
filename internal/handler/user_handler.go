package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/vladgrskkh/onerep-auth/internal/application"
	"github.com/vladgrskkh/onerep-auth/internal/domain"
)

type userProfileService interface {
	GetProfile(
		ctx context.Context,
		userID,
		requesterID uuid.UUID,
	) (domain.UserProfile, error)
	UpdateProfile(
		ctx context.Context,
		userID uuid.UUID,
		input application.UpdateProfileInput,
	) (domain.UserProfile, error)
}

type UserHandler struct {
	svc userProfileService
}

func NewUserHandler(svc userProfileService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Get("/v1/users/{id}", h.getProfile)
	r.Patch("/v1/users/{id}", h.updateProfile)
}

func (h *UserHandler) getProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	requesterID := userIDFromContext(r.Context())

	profile, err := h.svc.GetProfile(r.Context(), userID, requesterID)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, profile)
}

func (h *UserHandler) updateProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	requesterID := userIDFromContext(r.Context())
	if userID != requesterID {
		writeError(w, http.StatusForbidden, "cannot update another user's profile")
		return
	}

	var input application.UpdateProfileInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	profile, err := h.svc.UpdateProfile(r.Context(), userID, input)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, profile)
}
