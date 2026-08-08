package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/vladgrskkh/onerep-auth/internal/auth"
	"github.com/vladgrskkh/onerep-auth/internal/handler"
)

type userProfileService interface {
	GetProfile(
		ctx context.Context,
		userID,
		requesterID uuid.UUID,
	) (auth.UserProfile, error)
	UpdateProfile(
		ctx context.Context,
		userID uuid.UUID,
		input UpdateProfileInput,
	) (auth.UserProfile, error)
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
		handler.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	requesterID := handler.UserIDFromContext(r.Context())

	profile, err := h.svc.GetProfile(r.Context(), userID, requesterID)
	if err != nil {
		handler.WriteDomainError(w, err)
		return
	}

	handler.WriteJSON(w, http.StatusOK, profile)
}

func (h *UserHandler) updateProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	requesterID := handler.UserIDFromContext(r.Context())
	if userID != requesterID {
		handler.WriteError(w, http.StatusForbidden, "cannot update another user's profile")
		return
	}

	var input UpdateProfileInput
	if decErr := json.NewDecoder(r.Body).Decode(&input); decErr != nil {
		handler.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	profile, err := h.svc.UpdateProfile(r.Context(), userID, input)
	if err != nil {
		handler.WriteDomainError(w, err)
		return
	}

	handler.WriteJSON(w, http.StatusOK, profile)
}
