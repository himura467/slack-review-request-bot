package config

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/himura467/slack-review-request-bot/internal/domain/model"
)

var (
	OAuthToken    = ""
	SigningSecret = ""
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

	reviewerMapBytes, err := os.ReadFile("reviewer_map.json")
	if err != nil {
		slog.Error("failed to read reviewer map file", "error", err)
	}
	if err := json.Unmarshal(reviewerMapBytes, &reviewerMap); err != nil {
		slog.Error("failed to parse reviewer map config", "error", err)
	}

	return &SlackConfig{
		OAuthToken:    model.OAuthToken(token),
		SigningSecret: model.SigningSecret(secret),
		ReviewerMap:   reviewerMap,
	}
}
