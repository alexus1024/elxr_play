package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/caarlos0/env/v11"
)

type App struct {
	config     Config
	logger     *slog.Logger
	counter    int
	counterMu  sync.Mutex
	httpServer *http.Server
}

// @title Simple Go API
// @version 1.0
// @description A simple API with counter and fake LLM endpoints
// @BasePath /api
func New(ctx context.Context, cfg Config) (*App, error) {
	level := slog.LevelInfo
	if cfg.LogLevel == "debug" {
		level = slog.LevelDebug
	}

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(logHandler)

	app := &App{
		config: cfg,
		logger: logger,
	}

	logger.InfoContext(ctx, "app initialized", slog.Int("port", cfg.Port))
	return app, nil
}

func (a *App) Start(ctx context.Context) error {
	a.setupRoutes()
	a.logger.InfoContext(ctx, "starting server", slog.String("addr", a.httpServer.Addr))

	if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.logger.ErrorContext(ctx, "server error", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	a.logger.InfoContext(ctx, "stopping server")
	return a.httpServer.Shutdown(ctx)
}

func main() {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse config: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	app, err := New(ctx, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create app: %v\n", err)
		os.Exit(1)
	}

	if err := app.Start(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start server: %v\n", err)
		os.Exit(1)
	}
}
