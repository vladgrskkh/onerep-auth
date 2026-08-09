package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	jwtsvc "github.com/vladgrskkh/onerep-auth/internal/infrastructure/auth/jwt"

	"github.com/vladgrskkh/onerep-auth/internal/handler"
	"github.com/vladgrskkh/onerep-auth/internal/handler/auth/dto"
)

type authService interface {
	Register(ctx context.Context, email, password, displayName string) (jwtsvc.TokenPair, error)
	Login(ctx context.Context, email, password string) (jwtsvc.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	Refresh(ctx context.Context, refreshToken string) (jwtsvc.TokenPair, error)
}

type AuthHandler struct {
	svc    authService
	logger *slog.Logger
}

func NewAuthHandler(svc authService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{svc: svc, logger: logger}
}

// Register handles user registration.
//
// @Summary Register a new user
// @Description Create a new account with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration data"
// @Success 201 {object} jwtsvc.TokenPair
// @Failure 400 {object} handler.ErrorResponse
// @Failure 409 {object} handler.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := handler.DecodeAndValidate(r, &req); err != nil {
		if errors.Is(err, handler.ErrValidationFailed) {
			handler.WriteError(w, h.logger, http.StatusBadRequest, handler.ValidationErrorDetail())
			return
		}
		handler.WriteError(w, h.logger, http.StatusBadRequest, invalidRequestBodyDetail())
		return
	}

	pair, err := h.svc.Register(r.Context(), req.Email, req.Password, req.DisplayName)
	if err != nil {
		status, detail := mapError(err)
		handler.WriteError(w, h.logger, status, detail)
		return
	}

	handler.WriteJSON(w, h.logger, http.StatusCreated, pair)
}

// Login handles user authentication.
//
// @Summary Login
// @Description Authenticate with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login data"
// @Success 200 {object} jwtsvc.TokenPair
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := handler.DecodeAndValidate(r, &req); err != nil {
		if errors.Is(err, handler.ErrValidationFailed) {
			handler.WriteError(w, h.logger, http.StatusBadRequest, handler.ValidationErrorDetail())
			return
		}
		handler.WriteError(w, h.logger, http.StatusBadRequest, invalidRequestBodyDetail())
		return
	}

	pair, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		status, detail := mapError(err)
		handler.WriteError(w, h.logger, status, detail)
		return
	}

	handler.WriteJSON(w, h.logger, http.StatusOK, pair)
}

// Logout invalidates the refresh token.
//
// @Summary Logout
// @Description Invalidate refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LogoutRequest true "Logout data"
// @Success 204
// @Failure 400 {object} handler.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req dto.LogoutRequest
	if err := handler.DecodeAndValidate(r, &req); err != nil {
		if errors.Is(err, handler.ErrValidationFailed) {
			handler.WriteError(w, h.logger, http.StatusBadRequest, handler.ValidationErrorDetail())
			return
		}
		handler.WriteError(w, h.logger, http.StatusBadRequest, invalidRequestBodyDetail())
		return
	}

	if err := h.svc.Logout(r.Context(), req.RefreshToken); err != nil {
		status, detail := mapError(err)
		handler.WriteError(w, h.logger, status, detail)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Refresh rotates the refresh token.
//
// @Summary Refresh tokens
// @Description Get a new token pair using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshRequest true "Refresh data"
// @Success 200 {object} jwtsvc.TokenPair
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if err := handler.DecodeAndValidate(r, &req); err != nil {
		if errors.Is(err, handler.ErrValidationFailed) {
			handler.WriteError(w, h.logger, http.StatusBadRequest, handler.ValidationErrorDetail())
			return
		}
		handler.WriteError(w, h.logger, http.StatusBadRequest, invalidRequestBodyDetail())
		return
	}

	pair, err := h.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		status, detail := mapError(err)
		handler.WriteError(w, h.logger, status, detail)
		return
	}

	handler.WriteJSON(w, h.logger, http.StatusOK, pair)
}
