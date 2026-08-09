// @title OneRep Auth API
// @version 1.0
// @description Authentication service for OneRep gym training app
// @host localhost:8080
// @BasePath /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/vladgrskkh/onerep-auth/internal/application"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	app, err := application.New(application.WithLogger(logger))
	if err != nil {
		logger.Error("failed to create app", "error", err)
		return err
	}

	if err := app.Run(ctx); err != nil {
		logger.Error("server error", "error", err)
		return err
	}

	logger.Info("stopped")
	return nil
}
