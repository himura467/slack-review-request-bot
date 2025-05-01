//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/himura467/slack-review-request-bot/internal/config"
	"github.com/himura467/slack-review-request-bot/internal/domain/model"
	"github.com/himura467/slack-review-request-bot/internal/infrastructure"
	"github.com/himura467/slack-review-request-bot/internal/usecase"
)

func provideOAuthToken(cfg *config.SlackConfig) model.OAuthToken {
	return cfg.OAuthToken
}

func provideSigningSecret(cfg *config.SlackConfig) model.SigningSecret {
	return cfg.SigningSecret
}

func initializeApp() *app {
	wire.Build(
		config.NewSlackConfig,
		provideOAuthToken,
		provideSigningSecret,
		infrastructure.Set,
		usecase.Set,
		newApp,
	)
	return &app{}
}
