package wowdata

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type App struct {
	bnet *bnetClient
	// s3 client
	// db client
}

func NewApp(ctx context.Context, awsProfile string) (*App, error) {
	var app = App{}
	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(awsProfile))
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}
	logger(ctx).DebugContext(ctx, "aws config", "region", awsCfg.Region, "credentials", awsCfg.Credentials)

	bnetOAuth, err := getBnetOAuthClient(ctx, awsCfg)
	if err != nil {
		return nil, fmt.Errorf("get bnet oauth client: %w", err)
	}
	app.bnet, err = newBnet(ctx, bnetOAuth)
	if err != nil {
		return nil, fmt.Errorf("create bnet api client: %w", err)
	}
	return &app, nil
}

type oauthClient struct {
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
}

func getBnetOAuthClient(ctx context.Context, cfg aws.Config) (oauthClient, error) {
	client := secretsmanager.NewFromConfig(cfg)
	res, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String("wowdata/bnet/oauth"),
	})
	if err != nil {
		return oauthClient{}, fmt.Errorf("get secret value: %w", err)
	}
	var o oauthClient
	if err := json.Unmarshal([]byte(*res.SecretString), &o); err != nil {
		return oauthClient{}, fmt.Errorf("deserialize secret: %w", err)
	} else if o.ClientID == "" || o.ClientSecret == "" {
		return oauthClient{}, fmt.Errorf("unexpected secret format")
	}
	return o, nil
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
