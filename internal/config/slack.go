package config

import (
	"strings"

	"github.com/himura467/slack-review-request-bot/internal/domain/model"
)

var (
	OAuthToken    = ""
	SigningSecret = ""
	ReviewerIDs   = ""
)

type SlackConfig struct {
	OAuthToken    model.OAuthToken
	SigningSecret model.SigningSecret
	ReviewerIDs   model.ReviewerIDs
}

func NewSlackConfig() *SlackConfig {
	token := OAuthToken
	secret := SigningSecret
	reviewerIDs := strings.Split(ReviewerIDs, ",")

	return &SlackConfig{
		OAuthToken:    model.OAuthToken(token),
		SigningSecret: model.SigningSecret(secret),
		ReviewerIDs:   model.ReviewerIDs(reviewerIDs),
	}
}
