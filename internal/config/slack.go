package config

import (
	"encoding/json"
	"log/slog"

	"github.com/himura467/slack-review-request-bot/internal/domain/model"
)

var (
	OAuthToken    = ""
	SigningSecret = ""
	ReviewerMap   = ""
)

type SlackConfig struct {
	OAuthToken    model.OAuthToken
	SigningSecret model.SigningSecret
	ReviewerMap   model.ReviewerMap
}

func NewSlackConfig() *SlackConfig {
	token := OAuthToken
	secret := SigningSecret
	reviewerMap := make(model.ReviewerMap)

	if err := json.Unmarshal([]byte(ReviewerMap), &reviewerMap); err != nil {
		slog.Error("failed to parse reviewer map config", "error", err)
	}

	return &SlackConfig{
		OAuthToken:    model.OAuthToken(token),
		SigningSecret: model.SigningSecret(secret),
		ReviewerMap:   reviewerMap,
	}
}
