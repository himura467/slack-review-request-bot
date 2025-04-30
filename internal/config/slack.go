package config

import (
	"os"

	"github.com/himura467/slack-review-request-bot/internal/domain/model"
)

type SlackConfig struct {
	OAuthToken    model.OAuthToken
	SigningSecret model.SigningSecret
}

func NewSlackConfig() *SlackConfig {
	return &SlackConfig{
		OAuthToken:    model.OAuthToken(os.Getenv("SLACK_OAUTH_TOKEN")),
		SigningSecret: model.SigningSecret(os.Getenv("SLACK_SIGNING_SECRET")),
	}
}
