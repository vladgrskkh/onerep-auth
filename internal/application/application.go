package application

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	redislib "github.com/redis/go-redis/v9"

	"github.com/vladgrskkh/onerep-auth/internal/auth"
	"github.com/vladgrskkh/onerep-auth/internal/config"
	"github.com/vladgrskkh/onerep-auth/internal/handler"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/crypto"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/jwt"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/postgres"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/redis"
	"github.com/vladgrskkh/onerep-auth/internal/user"
)

type App struct {
	logger *slog.Logger
	cfg    *config.Config
	server *http.Server
	pool   *pgxpool.Pool
	redis  *redislib.Client
}

type Option func(*App)

func WithLogger(logger *slog.Logger) Option {
	return func(a *App) {
		a.logger = logger
	}
}

func WithDBPool(pool *pgxpool.Pool) Option {
	return func(a *App) {
		a.pool = pool
	}
}

func New(opts ...Option) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	app := &App{logger: logger, cfg: cfg}

	for _, opt := range opts {
		opt(app)
	}

	return app, nil
}

func (a *App) Run(ctx context.Context) error {
	if a.pool == nil {
		pool, err := pgxpool.New(ctx, a.cfg.DatabaseURL)
		if err != nil {
			return err
		}
		a.pool = pool
	}

	redisOpts, err := redislib.ParseURL(a.cfg.RedisURL)
	if err != nil {
		return err
	}
	redisClient := redislib.NewClient(redisOpts)
	a.redis = redisClient

	userRepo := postgres.NewUserRepo(a.pool)
	hasher := crypto.NewPasswordHasher()

	tm, err := jwt.NewTokenManager(a.cfg.JWTPrivateKeyPEM, a.cfg.JWTTokenTTL)
	if err != nil {
		return err
	}
	tokenStore := redis.NewTokenStore(redisClient, a.cfg.RefreshTokenTTL)

	authSvc := auth.NewAuthService(userRepo, *hasher, *tm, tokenStore)
	userSvc := user.NewUserService(userRepo, userRepo)

	r := chi.NewRouter()

	authHandler := auth.NewAuthHandler(authSvc)
	userHandler := user.NewUserHandler(userSvc)
	jwksHandler := handler.JWKSHandler(tm.PublicKey)

	RegisterRoutes(r, authHandler, userHandler, jwksHandler, tm, a.logger)

	a.server = &http.Server{
		Addr:              ":" + a.cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: time.Second * 5, //nolint:mnd // server read header timeout
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

func (a *App) ShutdownTimeout() time.Duration {
	return a.cfg.ShutdownTimeout
}

func (a *App) Shutdown(ctx context.Context) error {
	if a.server != nil {
		return a.server.Shutdown(ctx)
	}
	return nil
}
