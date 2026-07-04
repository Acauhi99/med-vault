package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Acauhi99/med-vault/internal/server"
	"github.com/Acauhi99/med-vault/internal/shared/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config failed", "error", err)
		os.Exit(1)
	}

	app, err := server.New(context.Background(), cfg, logger)
	if err != nil {
		slog.Error("server init failed", "error", err)
		os.Exit(1)
	}

	go func() {
		slog.Info("server starting", "addr", app.Server.Addr)
		if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.Close()
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)

	if err := app.Server.Shutdown(ctx); err != nil {
		cancel()
		app.Close()
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}
	cancel()
	app.Close()

	slog.Info("server stopped")
}
