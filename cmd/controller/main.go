// Command controller starts the Platform-Infra control plane.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ikwukao/platform-infra/internal/api"
	"github.com/ikwukao/platform-infra/internal/config"
	"github.com/ikwukao/platform-infra/internal/controller"
	"github.com/ikwukao/platform-infra/internal/deployments"
	"github.com/ikwukao/platform-infra/internal/projects"
	"github.com/ikwukao/platform-infra/internal/services"
	"github.com/ikwukao/platform-infra/internal/storage"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	))

	logger.Info(
		"starting platform-infra controller",
		"port",
		cfg.ServerPort,
	)

	ctx := context.Background()

	db, err := storage.NewPostgres(
		ctx,
		cfg.DatabaseURL,
	)
	if err != nil {
		logger.Error(
			"database connection failed",
			"error",
			err,
		)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Migrate(ctx); err != nil {
		logger.Error(
			"database migration failed",
			"error",
			err,
		)
		os.Exit(1)
	}

	projectRepository := projects.NewPostgresRepository(db)
	serviceRepository := services.NewPostgresRepository(db)
	deploymentRepository := deployments.NewPostgresRepository(db)

	deploymentController := controller.New(deploymentRepository)

	server := &http.Server{
		Addr: ":" + cfg.ServerPort,
		Handler: api.NewServer(
			projectRepository,
			serviceRepository,
			deploymentRepository,
		).Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	signalCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		logger.Info(
			"HTTP server listening",
			"addr",
			server.Addr,
		)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			logger.Error(
				"HTTP server failed",
				"error",
				err,
			)
			stop()
		}
	}()

	go func() {
		if err := deploymentController.Run(signalCtx, 5*time.Second); err != nil &&
			!errors.Is(err, context.Canceled) {
			logger.Error(
				"deployment controller failed",
				"error",
				err,
			)
			stop()
		}
	}()

	<-signalCtx.Done()

	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(
			"graceful shutdown failed",
			"error",
			err,
		)
		return
	}

	logger.Info("controller stopped")
}
