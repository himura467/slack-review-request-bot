package config

import (
	"os"
	"strings"

	"github.com/himura467/slack-review-request-bot/internal/domain/model"
)

type SlackConfig struct {
	OAuthToken    model.OAuthToken
	SigningSecret model.SigningSecret
	ReviewerIDs   model.ReviewerIDs
}

func NewSlackConfig() *SlackConfig {
	var reviewerIDs []string
	if ids := os.Getenv("SLACK_REVIEWER_IDS"); ids != "" {
		reviewerIDs = strings.Split(ids, ",")
	}

	return &SlackConfig{
		OAuthToken:    model.OAuthToken(os.Getenv("SLACK_OAUTH_TOKEN")),
		SigningSecret: model.SigningSecret(os.Getenv("SLACK_SIGNING_SECRET")),
		ReviewerIDs:   reviewerIDs,
	}
}
