package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"wowdata/internal/wowdata"
)

func main() {
	ctx, unregisterSignal := signal.NotifyContext(context.Background(), os.Interrupt)
	defer unregisterSignal()

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(logHandler)
	ctx = wowdata.LoggerContext(ctx, logger)

	app, err := wowdata.NewApp(ctx, "wowdata")
	if err != nil {
		logger.ErrorContext(ctx, "init failed", "err", err)
		os.Exit(1)
	}
	if err = app.Run(ctx); err != nil {
		logger.ErrorContext(ctx, "run failed", "err", err)
		os.Exit(1)
	}
}
