package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	decibelindexer "github.com/cresendoo/decidash-backend/internal/application/decibel-indexer"
	"github.com/cresendoo/decidash-backend/pkg/xlogger"
	"github.com/phsym/console-slog"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	var cfg decibelindexer.Config
	if err := cfg.Load(); err != nil {
		panic(err)
	}

	logger := xlogger.Build(
		ctx,
		[]slog.Handler{console.NewHandler(os.Stdout, &console.HandlerOptions{Level: cfg.Log.Level})},
		xlogger.WithRelease("v0.0.1"),
		xlogger.WithSentryDSN(cfg.SentryDSN),
		xlogger.WithNamespace("decibel-indexer"),
		xlogger.WithLogLevel(cfg.Log.Level),
		xlogger.WithLogFile(cfg.Log.File),
		xlogger.WithLogFormat(cfg.Log.Format),
	)
	slog.SetDefault(logger)

	app, err := decibelindexer.NewApplication(ctx, logger, &cfg)
	if err != nil {
		slog.Error("failed to create application", "error", err)
		return
	}

	if err := app.Start(); err != nil {
		slog.Error("failed to start application", "error", err)
		return
	}
	slog.Info("application started")

	<-ctx.Done()

	if err := app.Close(); err != nil {
		slog.Error("failed to close application", "error", err)
	}
	slog.Info("application closed")
}
