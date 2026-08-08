package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/vladgrskkh/onerep-auth/internal/handler/middleware"
)

func RegisterRoutes(
	r chi.Router,
	authHandler *AuthHandler,
	userHandler *UserHandler,
	jwksHandler http.HandlerFunc,
	tokenValidator middleware.JWTValidator,
	logger *slog.Logger,
) {
	r.Use(middleware.Logging(logger))
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)

	r.Get("/health", healthHandler)
	r.Get("/.well-known/jwks.json", jwksHandler)

	authHandler.RegisterRoutes(r)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate(tokenValidator))
		userHandler.RegisterRoutes(r)
	})
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
