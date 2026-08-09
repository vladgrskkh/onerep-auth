package user

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	authdomain "github.com/vladgrskkh/onerep-auth/internal/domain/auth"
	userdomain "github.com/vladgrskkh/onerep-auth/internal/domain/user"
	"github.com/vladgrskkh/onerep-auth/internal/handler"
	"github.com/vladgrskkh/onerep-auth/internal/handler/user/dto"
)

const (
	errCodeInvalidUserID = "INVALID_USER_ID"
	errMsgInvalidUserID  = "invalid user id"
	errUserInvalidUserID = "The user ID is invalid"

	errCodeForbidden = "FORBIDDEN"
	errMsgForbidden  = "cannot update another user's profile"
	errUserForbidden = "You can only update your own profile"
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

// GetProfile returns the user's profile.
//
// @Summary Get user profile
// @Description Get a user's public profile (email visible only to owner)
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} authdomain.UserProfile
// @Failure 400 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, parseErr := uuid.Parse(chi.URLParam(r, "id"))
	if parseErr != nil {
		handler.WriteError(w, http.StatusBadRequest, handler.ErrorDetail{
			Code:        errCodeInvalidUserID,
			Message:     errMsgInvalidUserID,
			UserMessage: errUserInvalidUserID,
			Type:        handler.ErrorTypeUser,
		})
		return
	}

	requesterID := handler.UserIDFromContext(r.Context())

	profile, err := h.svc.GetProfile(r.Context(), userID, requesterID)
	if err != nil {
		status, detail := mapError(err)
		handler.WriteError(w, status, detail)
		return
	}

	handler.WriteJSON(w, http.StatusOK, profile)
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
// @Success 200 {object} authdomain.UserProfile
// @Failure 400 {object} handler.ErrorResponse
// @Failure 403 {object} handler.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [patch]
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, parseErr := uuid.Parse(chi.URLParam(r, "id"))
	if parseErr != nil {
		handler.WriteError(w, http.StatusBadRequest, handler.ErrorDetail{
			Code:        errCodeInvalidUserID,
			Message:     errMsgInvalidUserID,
			UserMessage: errUserInvalidUserID,
			Type:        handler.ErrorTypeUser,
		})
		return
	}

	requesterID := handler.UserIDFromContext(r.Context())
	if userID != requesterID {
		handler.WriteError(w, http.StatusForbidden, handler.ErrorDetail{
			Code:        errCodeForbidden,
			Message:     errMsgForbidden,
			UserMessage: errUserForbidden,
			Type:        handler.ErrorTypeUser,
		})
		return
	}

	var req dto.UpdateProfileRequest
	if decErr := handler.DecodeJSON(r, &req); decErr != nil {
		handler.WriteError(w, http.StatusBadRequest, handler.ErrorDetail{
			Code:        "INVALID_REQUEST_BODY",
			Message:     "invalid request body",
			UserMessage: "The request body is invalid",
			Type:        handler.ErrorTypeUser,
		})
		return
	}

	input := userdomain.UpdateProfileInput{
		DisplayName: req.DisplayName,
		BirthDate:   nil, // TODO: parse from req.BirthDate
	}
	if req.Gender != nil {
		g := authdomain.Gender(*req.Gender)
		input.Gender = &g
	}

	profile, err := h.svc.UpdateProfile(r.Context(), userID, input)
	if err != nil {
		status, detail := mapError(err)
		handler.WriteError(w, status, detail)
		return
	}

	handler.WriteJSON(w, http.StatusOK, profile)
}
