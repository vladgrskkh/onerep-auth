package application

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

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
	logger        *slog.Logger
	cfg           *config.Config
	server        *http.Server
	pool          *pgxpool.Pool
	redis         *redislib.Client
	tokenManager  *jwt.TokenManager
	authHandler   *authhandler.AuthHandler
	userHandler   *userhandler.UserHandler
	jwksHandler   http.HandlerFunc
	shutdownTimer time.Duration
}

type Option func(*Application)

func WithLogger(logger *slog.Logger) Option {
	return func(a *Application) {
		a.logger = logger
	}
}

func WithDBPool(pool *pgxpool.Pool) Option {
	return func(a *Application) {
		a.pool = pool
	}
}

func WithRedisClient(client *redislib.Client) Option {
	return func(a *Application) {
		a.redis = client
	}
}

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(a *Application) {
		a.shutdownTimer = timeout
	}
}

func New(opts ...Option) (*Application, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	app := &Application{
		logger:        logger,
		cfg:           cfg,
		shutdownTimer: cfg.ShutdownTimeout,
	}

	for _, opt := range opts {
		opt(app)
	}

	return app, nil
}

func (app *Application) connectDependencies(ctx context.Context) error {
	if app.pool == nil {
		pool, err := pgxpool.New(ctx, app.cfg.DatabaseURL)
		if err != nil {
			return err
		}
		app.pool = pool
	}

	if app.redis == nil {
		redisOpts, err := redislib.ParseURL(app.cfg.RedisURL)
		if err != nil {
			return err
		}
		app.redis = redislib.NewClient(redisOpts)
	}

	return nil
}

func (app *Application) Run(ctx context.Context) error {
	if err := app.connectDependencies(ctx); err != nil {
		return err
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
	userSvc := userservice.NewUserService(userRepo, userRepo)

	app.authHandler = authhandler.NewAuthHandler(authSvc)
	app.userHandler = userhandler.NewUserHandler(userSvc)
	app.jwksHandler = handler.JWKSHandler(tm.PublicKey)

	app.server = &http.Server{
		Addr:              ":" + app.cfg.Port,
		Handler:           app.RegisterRoutes(tm),
		ReadHeaderTimeout: time.Second * 5, //nolint:mnd // server read header timeout
	}

	go func() {
		app.logger.Info("starting auth service", "port", app.cfg.Port)
		if listenErr := app.server.ListenAndServe(); listenErr != nil && listenErr != http.ErrServerClosed {
			app.logger.Error("server error", "error", listenErr)
			os.Exit(1)
		}
	}()

	return nil
}

func (app *Application) Shutdown(ctx context.Context) error {
	if app.server != nil {
		return app.server.Shutdown(ctx)
	}
	return nil
}

func (app *Application) ShutdownTimeout() time.Duration {
	return app.shutdownTimer
}
