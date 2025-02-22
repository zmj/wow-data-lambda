package wowdata

import (
	"context"
	"log/slog"
)

type App struct {
	// bnet api client
	// s3 client
	// db client
}

func NewApp(ctx context.Context, awsProfile string) (*App, error) {
	return &App{}, nil
}

type loggerContextKey struct{}

func LoggerContext(parent context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(parent, loggerContextKey{}, logger)
}

func logger(ctx context.Context) *slog.Logger {
	val := ctx.Value(loggerContextKey{})
	if val == nil {
		panic("logger not set")
	}
	logger, ok := val.(*slog.Logger)
	if !ok {
		panic("unexpected logger type")
	}
	return logger
}
