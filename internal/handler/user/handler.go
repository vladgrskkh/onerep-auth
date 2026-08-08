package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	authdomain "github.com/vladgrskkh/onerep-auth/internal/domain/auth"
	userdomain "github.com/vladgrskkh/onerep-auth/internal/domain/user"
	"github.com/vladgrskkh/onerep-auth/internal/handler"
)

type userProfileService interface {
	GetProfile(
		ctx context.Context,
		userID,
		requesterID uuid.UUID,
	) (authdomain.UserProfile, error)
	UpdateProfile(
		ctx context.Context,
		userID uuid.UUID,
		input userdomain.UpdateProfileInput,
	) (authdomain.UserProfile, error)
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

// @Summary Get user profile
// @Description Get a user's public profile (email visible only to owner)
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} authdomain.UserProfile
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /users/{id} [get]
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

// @Summary Update user profile
// @Description Update the authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body userdomain.UpdateProfileInput true "Profile data"
// @Success 200 {object} authdomain.UserProfile
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Security BearerAuth
// @Router /users/{id} [patch]
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

	var input userdomain.UpdateProfileInput
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
