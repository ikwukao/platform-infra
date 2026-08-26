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
	"github.com/ikwukao/platform-infra/internal/projects"
	"github.com/ikwukao/platform-infra/internal/storage"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	logger.Info("starting platform-infra controller",
		"port", cfg.ServerPort,
	)

	db, err := storage.NewPostgres(
		context.Background(),
		cfg.DatabaseURL,
	)
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	projectRepository := projects.NewPostgresRepository(db)

	server := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           api.NewServer(projectRepository).Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		logger.Info("HTTP server listening", "addr", server.Addr)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server failed", "error", err)
			stop()
		}
	}()

	<-ctx.Done()

	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		return
	}

	logger.Info("controller stopped")
}
