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

func provideReviewerIDs(cfg *config.SlackConfig) model.ReviewerIDs {
	return cfg.ReviewerIDs
}

func initializeApp() *app {
	wire.Build(
		config.NewSlackConfig,
		rest.Set,
		provideOAuthToken,
		provideSigningSecret,
		provideReviewerIDs,
		newApp,
	)
	return &app{}
}
