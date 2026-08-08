package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/vladgrskkh/onerep-auth/internal/handler"
	jwtsvc "github.com/vladgrskkh/onerep-auth/internal/infrastructure/jwt"
)

type authService interface {
	Register(ctx context.Context, email, password, displayName string) (jwtsvc.TokenPair, error)
	Login(ctx context.Context, email, password string) (jwtsvc.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	Refresh(ctx context.Context, refreshToken string) (jwtsvc.TokenPair, error)
}

type AuthHandler struct {
	svc authService
}

func NewAuthHandler(svc authService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Post("/v1/auth/register", h.register)
	r.Post("/v1/auth/login", h.login)
	r.Post("/v1/auth/logout", h.logout)
	r.Post("/v1/auth/refresh", h.refresh)
}

func (h *AuthHandler) register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		DisplayName string `json:"display_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	pair, err := h.svc.Register(r.Context(), req.Email, req.Password, req.DisplayName)
	if err != nil {
		handler.WriteDomainError(w, err)
		return
	}

	handler.WriteJSON(w, http.StatusCreated, pair)
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	pair, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		handler.WriteDomainError(w, err)
		return
	}

	handler.WriteJSON(w, http.StatusOK, pair)
}

func (h *AuthHandler) logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.svc.Logout(r.Context(), req.RefreshToken); err != nil {
		handler.WriteDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	pair, err := h.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		handler.WriteDomainError(w, err)
		return
	}

	handler.WriteJSON(w, http.StatusOK, pair)
}
