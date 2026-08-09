package application

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	redislib "github.com/redis/go-redis/v9"

	"github.com/vladgrskkh/onerep-auth/internal/config"
	"github.com/vladgrskkh/onerep-auth/internal/handler"
	authhandler "github.com/vladgrskkh/onerep-auth/internal/handler/auth"
	userhandler "github.com/vladgrskkh/onerep-auth/internal/handler/user"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/auth/crypto"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/auth/jwt"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/auth/postgres"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/auth/redis"
	authservice "github.com/vladgrskkh/onerep-auth/internal/service/auth"
	userservice "github.com/vladgrskkh/onerep-auth/internal/service/user"
)

type Application struct {
	logger       *slog.Logger
	cfg          *config.Config
	server       *http.Server
	pool         *pgxpool.Pool
	redis        *redislib.Client
	tokenManager *jwt.TokenManager
	authHandler  *authhandler.AuthHandler
	userHandler  *userhandler.UserHandler
	jwksHandler  http.HandlerFunc
}

type Option func(*Application) error

func WithLogger(logger *slog.Logger) Option {
	return func(a *Application) error {
		a.logger = logger
		return nil
	}
}

// WithDatabase connects to PostgreSQL using the configured URL.
func WithDatabase(ctx context.Context) Option {
	return func(a *Application) error {
		pool, err := pgxpool.New(ctx, a.cfg.DatabaseURL)
		if err != nil {
			return fmt.Errorf("connect to postgres: %w", err)
		}
		a.pool = pool
		return nil
	}
}

// WithRedis connects to Redis using the configured URL.
func WithRedis(_ context.Context) Option {
	return func(a *Application) error {
		opts, err := redislib.ParseURL(a.cfg.RedisURL)
		if err != nil {
			return fmt.Errorf("parse redis url: %w", err)
		}
		a.redis = redislib.NewClient(opts)
		return nil
	}
}

func New(opts ...Option) (*Application, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	app := &Application{
		logger: logger,
		cfg:    cfg,
	}

	for _, opt := range opts {
		if err := opt(app); err != nil {
			return nil, err
		}
	}

	return app, nil
}

func (app *Application) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()

	if app.pool == nil {
		return errors.New("database pool not initialized: use application.WithDatabase")
	}
	if app.redis == nil {
		return errors.New("redis client not initialized: use application.WithRedis")
	}

	tm, err := jwt.NewTokenManager(app.cfg.JWTPrivateKeyPEM, app.cfg.JWTTokenTTL)
	if err != nil {
		return err
	}
	app.tokenManager = tm

	userRepo := postgres.NewUserRepo(app.pool)
	hasher := crypto.NewPasswordHasher()
	tokenStore := redis.NewTokenStore(app.redis, app.cfg.RefreshTokenTTL)

	authSvc := authservice.NewAuthService(userRepo, *hasher, *tm, tokenStore)
	userSvc := userservice.NewUserService(userRepo)

	app.authHandler = authhandler.NewAuthHandler(authSvc, app.logger)
	app.userHandler = userhandler.NewUserHandler(userSvc, app.logger)
	app.jwksHandler = handler.JWKSHandler(tm.PublicKey, app.logger)

	app.server = &http.Server{
		Addr:              ":" + app.cfg.Port,
		Handler:           app.RegisterRoutes(tm),
		ReadHeaderTimeout: app.cfg.ReadHeaderTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		app.logger.Info("starting auth service", "port", app.cfg.Port)
		errCh <- app.server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		app.logger.Info("shutting down...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), app.cfg.ShutdownTimeout)
		defer shutdownCancel()

		if err := app.server.Shutdown(shutdownCtx); err != nil {
			return err
		}
		return nil
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
}

func (app *Application) Shutdown(ctx context.Context) error {
	if app.server != nil {
		return app.server.Shutdown(ctx)
	}
	return nil
}
