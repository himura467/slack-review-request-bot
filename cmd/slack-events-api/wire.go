//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/himura467/slack-review-request-bot/internal/config"
	"github.com/himura467/slack-review-request-bot/internal/domain/model"
	"github.com/himura467/slack-review-request-bot/internal/interface/rest"
)

func provideOAuthToken(cfg *config.SlackConfig) model.OAuthToken {
	return cfg.OAuthToken
}

func provideSigningSecret(cfg *config.SlackConfig) model.SigningSecret {
	return cfg.SigningSecret
}

func provideReviewerMap(cfg *config.SlackConfig) model.ReviewerMap {
	return cfg.ReviewerMap
}

func initializeApp() *app {
	wire.Build(
		config.NewSlackConfig,
		rest.Set,
		provideOAuthToken,
		provideSigningSecret,
		provideReviewerMap,
		newApp,
	)
	return &app{}
}
