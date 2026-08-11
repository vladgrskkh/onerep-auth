package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	authdomain "github.com/vladgrskkh/onerep-auth/internal/domain/auth"
	jwtsvc "github.com/vladgrskkh/onerep-auth/internal/infrastructure/auth/jwt"
	authsvc "github.com/vladgrskkh/onerep-auth/internal/service/auth"

	"github.com/vladgrskkh/onerep-auth/internal/handler"
	"github.com/vladgrskkh/onerep-auth/internal/handler/auth/dto"
)

// AuthService is the auth use-case contract consumed by the handler.
type AuthService interface {
	Register(ctx context.Context, cmd authsvc.RegisterCommand) (jwtsvc.TokenPair, error)
	Login(ctx context.Context, cmd authsvc.LoginCommand) (jwtsvc.TokenPair, error)
	Logout(ctx context.Context, cmd authsvc.LogoutCommand) error
	Refresh(ctx context.Context, cmd authsvc.RefreshCommand) (jwtsvc.TokenPair, error)
}

type AuthHandler struct {
	svc    AuthService
	logger *slog.Logger
}

func NewAuthHandler(svc AuthService, logger *slog.Logger) *AuthHandler {
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
// @Failure 500 {object} handler.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := handler.DecodeAndValidate(r, &req); err != nil {
		if errors.Is(err, handler.ErrValidationFailed) {
			handler.WriteError(w, h.logger, http.StatusBadRequest, handler.ValidationErrorDetail(err))
			return
		}
		handler.WriteError(w, h.logger, http.StatusBadRequest, invalidRequestBodyDetail())
		return
	}

	pair, err := h.svc.Register(r.Context(), authsvc.RegisterCommand{
		Email:       req.Email,
		Password:    req.Password,
		DisplayName: req.DisplayName,
	})
	if err != nil {
		switch {
		case errors.Is(err, authdomain.ErrEmailAlreadyExists):
			handler.WriteError(w, h.logger, http.StatusConflict, emailAlreadyExistsDetail())
		case errors.Is(err, authdomain.ErrInvalidEmail):
			handler.WriteError(w, h.logger, http.StatusBadRequest, invalidEmailDetail())
		case errors.Is(err, authdomain.ErrInvalidPassword):
			handler.WriteError(w, h.logger, http.StatusBadRequest, invalidPasswordDetail())
		case errors.Is(err, authdomain.ErrInvalidDisplayName):
			handler.WriteError(w, h.logger, http.StatusBadRequest, invalidDisplayNameDetail())
		default:
			handler.WriteSystemError(w, h.logger, http.StatusInternalServerError, err)
		}
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
// @Failure 500 {object} handler.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := handler.DecodeAndValidate(r, &req); err != nil {
		if errors.Is(err, handler.ErrValidationFailed) {
			handler.WriteError(w, h.logger, http.StatusBadRequest, handler.ValidationErrorDetail(err))
			return
		}
		handler.WriteError(w, h.logger, http.StatusBadRequest, invalidRequestBodyDetail())
		return
	}

	pair, err := h.svc.Login(r.Context(), authsvc.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, authdomain.ErrInvalidCredentials):
			handler.WriteError(w, h.logger, http.StatusUnauthorized, invalidCredentialsDetail())
		default:
			handler.WriteSystemError(w, h.logger, http.StatusInternalServerError, err)
		}
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
// @Failure 500 {object} handler.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req dto.LogoutRequest
	if err := handler.DecodeAndValidate(r, &req); err != nil {
		if errors.Is(err, handler.ErrValidationFailed) {
			handler.WriteError(w, h.logger, http.StatusBadRequest, handler.ValidationErrorDetail(err))
			return
		}
		handler.WriteError(w, h.logger, http.StatusBadRequest, invalidRequestBodyDetail())
		return
	}

	if err := h.svc.Logout(r.Context(), authsvc.LogoutCommand{RefreshToken: req.RefreshToken}); err != nil {
		handler.WriteSystemError(w, h.logger, http.StatusInternalServerError, err)
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
// @Failure 500 {object} handler.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if err := handler.DecodeAndValidate(r, &req); err != nil {
		if errors.Is(err, handler.ErrValidationFailed) {
			handler.WriteError(w, h.logger, http.StatusBadRequest, handler.ValidationErrorDetail(err))
			return
		}
		handler.WriteError(w, h.logger, http.StatusBadRequest, invalidRequestBodyDetail())
		return
	}

	pair, err := h.svc.Refresh(r.Context(), authsvc.RefreshCommand{RefreshToken: req.RefreshToken})
	if err != nil {
		switch {
		case errors.Is(err, authdomain.ErrTokenNotFound):
			handler.WriteError(w, h.logger, http.StatusUnauthorized, tokenNotFoundDetail())
		default:
			handler.WriteSystemError(w, h.logger, http.StatusInternalServerError, err)
		}
		return
	}

	handler.WriteJSON(w, h.logger, http.StatusOK, pair)
}
