package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type App struct {
	config     Config
	logger     *slog.Logger
	counter    int
	counterMu  sync.Mutex
	httpServer *http.Server
	quit       chan struct{}
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
		quit:   make(chan struct{}),
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
	cfg, err := LoadConfig()
	if err != nil {
		slog.Error("Can not load config", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	app, err := New(ctx, cfg)
	if err != nil {
		slog.Error("failed to create app", "error", err)
		os.Exit(1)
	}

	go func() {
		if err := app.Start(ctx); err != nil {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
	case <-app.quit:
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, cfg.ShutdownTimeout)
	defer cancel()
	if err := app.Stop(shutdownCtx); err != nil {
		slog.Error("failed to stop server", "error", err)
	}

	app.logger.InfoContext(ctx, "server stopped gracefully")
}
