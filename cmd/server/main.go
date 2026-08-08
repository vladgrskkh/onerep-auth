package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	redislib "github.com/redis/go-redis/v9"

	"github.com/vladgrskkh/onerep-auth/internal/application"
	"github.com/vladgrskkh/onerep-auth/internal/config"
	"github.com/vladgrskkh/onerep-auth/internal/handler"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/crypto"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/jwt"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/postgres"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/redis"
)

type App struct {
	logger *slog.Logger
	cfg    *config.Config
	server *http.Server
	pool   *pgxpool.Pool
	redis  *redislib.Client
}

func NewApp(opts ...Option) *App {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	app := &App{logger: logger, cfg: cfg}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

func (a *App) Run(ctx context.Context) error {
	pool, err := pgxpool.New(ctx, a.cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	a.pool = pool

	redisOpts, err := redislib.ParseURL(a.cfg.RedisURL)
	if err != nil {
		return err
	}
	redisClient := redislib.NewClient(redisOpts)
	a.redis = redisClient

	userRepo := postgres.NewUserRepo(pool)
	hasher := crypto.NewPasswordHasher()

	tm, err := jwt.NewTokenManager(a.cfg.JWTPrivateKeyPEM, a.cfg.JWTTokenTTL)
	if err != nil {
		return err
	}
	tokenStore := redis.NewTokenStore(redisClient, a.cfg.RefreshTokenTTL)

	authSvc := application.NewAuthService(userRepo, *hasher, *tm, tokenStore)
	userSvc := application.NewUserService(userRepo, userRepo)

	r := chi.NewRouter()

	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	jwksHandler := handler.JWKSHandler(tm.PublicKey)

	handler.RegisterRoutes(r, authHandler, userHandler, jwksHandler, tm, a.logger)

	a.server = &http.Server{
		Addr:              ":" + a.cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		a.logger.Info("starting auth service", "port", a.cfg.Port)
		if listenErr := a.server.ListenAndServe(); listenErr != nil && listenErr != http.ErrServerClosed {
			a.logger.Error("server error", "error", listenErr)
			os.Exit(1)
		}	
	}()

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	if a.server != nil {
		return a.server.Shutdown(ctx)
	}
	return nil
}

type Option func(*App)

func WithLogger(logger *slog.Logger) Option {
	return func(a *App) {
		a.logger = logger
	}
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	app := NewApp()

	if err := app.Run(ctx); err != nil {
		app.logger.Error("failed to start", "error", err)
		return err
	}

	<-ctx.Done()
	app.logger.Info("shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		app.logger.Error("force shutdown", "error", err)
		return err
	}

	app.logger.Info("stopped")
	return nil
}
