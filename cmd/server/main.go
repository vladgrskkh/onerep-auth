package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	redislib "github.com/redis/go-redis/v9"

	"github.com/vladgrskkh/onerep-auth/internal/application"
	"github.com/vladgrskkh/onerep-auth/internal/handler"
	"github.com/vladgrskkh/onerep-auth/internal/handler/middleware"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/crypto"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/jwt"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/postgres"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/redis"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	databaseURL := os.Getenv("DATABASE_URL")
	redisURL := os.Getenv("REDIS_URL")
	port := envOrDefault("AUTH_PORT", "8080")
	jwtPrivateKeyPEM := os.Getenv("JWT_PRIVATE_KEY")
	jwtTTL := 15 * time.Minute
	refreshTTL := 7 * 24 * time.Hour

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		logger.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	redisOpts, err := redislib.ParseURL(redisURL)
	if err != nil {
		logger.Error("failed to parse redis URL", "error", err)
		os.Exit(1)
	}
	redisClient := redislib.NewClient(redisOpts)

	userRepo := postgres.NewUserRepo(pool)

	hasher := crypto.NewPasswordHasher()
	jwtSvc, err := jwt.NewService(jwtPrivateKeyPEM, jwtTTL)
	if err != nil {
		logger.Error("failed to create JWT service", "error", err)
		os.Exit(1)
	}
	tokenStore := redis.NewTokenStore(redisClient, refreshTTL)

	authSvc := application.NewAuthService(userRepo, userRepo, userRepo, *hasher, *jwtSvc, tokenStore)
	userSvc := application.NewUserService(userRepo, userRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logging(logger))
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)

	r.Get("/health", healthHandler)
	r.Get("/.well-known/jwks.json", handler.JWKSHandler(jwtSvc.PublicKey))

	authHandler := handler.NewAuthHandler(authSvc)
	authHandler.RegisterRoutes(r)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate(jwtSvc))
		userHandler := handler.NewUserHandler(userSvc)
		userHandler.RegisterRoutes(r)
	})

	logger.Info("starting auth service", "port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
