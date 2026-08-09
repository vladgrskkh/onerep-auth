package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/vladgrskkh/onerep-auth/docs" // swagger docs
	"github.com/vladgrskkh/onerep-auth/internal/handler"
	"github.com/vladgrskkh/onerep-auth/internal/handler/middleware"
)

func (app *Application) RegisterRoutes(tokenValidator middleware.JWTValidator) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logging(app.logger))
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)

	r.Get("/v1/health", healthHandler)
	r.Get("/v1/.well-known/jwks.json", app.jwksHandler)
	r.Get("/v1/docs/*", httpSwagger.WrapHandler)

	r.Post("/v1/auth/register", app.authHandler.Register)
	r.Post("/v1/auth/login", app.authHandler.Login)
	r.Post("/v1/auth/logout", app.authHandler.Logout)
	r.Post("/v1/auth/refresh", app.authHandler.Refresh)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate(tokenValidator))
		r.Get("/v1/users/{id}", app.userHandler.GetProfile)
		r.Patch("/v1/users/{id}", app.userHandler.UpdateProfile)
	})

	return r
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	handler.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
