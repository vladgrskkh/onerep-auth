package user

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	userdomain "github.com/vladgrskkh/onerep-auth/internal/domain/user"
	"github.com/vladgrskkh/onerep-auth/internal/handler"
	"github.com/vladgrskkh/onerep-auth/internal/handler/user/dto"
)

type userProfileService interface {
	GetProfile(ctx context.Context, userID, requesterID uuid.UUID) (userdomain.UserProfile, error)
	UpdateProfile(ctx context.Context, profile userdomain.UserProfile) (userdomain.UserProfile, error)
}

type UserHandler struct {
	svc    userProfileService
	logger *slog.Logger
}

func NewUserHandler(svc userProfileService, logger *slog.Logger) *UserHandler {
	return &UserHandler{svc: svc, logger: logger}
}

// GetProfile returns the user's profile.
//
// @Summary Get user profile
// @Description Get a user's public profile (email visible only to owner)
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserProfileResponse
// @Failure 400 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, parseErr := uuid.Parse(chi.URLParam(r, "id"))
	if parseErr != nil {
		handler.WriteError(w, h.logger, http.StatusBadRequest, invalidUserIDDetail())
		return
	}

	requesterID := handler.UserIDFromContext(r.Context())

	profile, err := h.svc.GetProfile(r.Context(), userID, requesterID)
	if err != nil {
		status, detail := mapError(err)
		handler.WriteError(w, h.logger, status, detail)
		return
	}

	handler.WriteJSON(w, h.logger, http.StatusOK, toProfileResponse(profile))
}

// UpdateProfile updates the authenticated user's profile.
//
// @Summary Update user profile
// @Description Update the authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body dto.UpdateProfileRequest true "Profile data"
// @Success 200 {object} dto.UserProfileResponse
// @Failure 400 {object} handler.ErrorResponse
// @Failure 403 {object} handler.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [patch]
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, parseErr := uuid.Parse(chi.URLParam(r, "id"))
	if parseErr != nil {
		handler.WriteError(w, h.logger, http.StatusBadRequest, invalidUserIDDetail())
		return
	}

	requesterID := handler.UserIDFromContext(r.Context())
	if userID != requesterID {
		handler.WriteError(w, h.logger, http.StatusForbidden, forbiddenDetail())
		return
	}

	var req dto.UpdateProfileRequest
	if err := handler.DecodeAndValidate(r, &req); err != nil {
		if errors.Is(err, handler.ErrValidationFailed) {
			handler.WriteError(w, h.logger, http.StatusBadRequest, handler.ValidationErrorDetail())
			return
		}
		handler.WriteError(w, h.logger, http.StatusBadRequest, invalidRequestBodyDetail())
		return
	}

	profile, detail := toUpdateProfile(req, userID)
	if detail.Code != "" {
		handler.WriteError(w, h.logger, http.StatusBadRequest, detail)
		return
	}

	updated, err := h.svc.UpdateProfile(r.Context(), profile)
	if err != nil {
		status, detail := mapError(err)
		handler.WriteError(w, h.logger, status, detail)
		return
	}

	handler.WriteJSON(w, h.logger, http.StatusOK, toProfileResponse(updated))
}
