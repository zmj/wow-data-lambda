package main

import (
	"context"
	"log/slog"
	"os"
	"wowdata/internal/wowdata"

	"github.com/aws/aws-lambda-go/lambda"
)

var logger *slog.Logger
var app *wowdata.App

func init() {
	logHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger = slog.New(logHandler)
	ctx := wowdata.LoggerContext(context.Background(), logger)

	var err error
	app, err = wowdata.NewApp(ctx, "")
	if err != nil {
		logger.ErrorContext(ctx, "init failed", "err", err)
		os.Exit(1)
	}
}

func main() {
	lambda.Start(func(ctx context.Context) {
		ctx = wowdata.LoggerContext(ctx, logger)
		if err := app.Run(ctx); err != nil {
			logger.ErrorContext(ctx, "run failed", "err", err)
		}
	})
}
