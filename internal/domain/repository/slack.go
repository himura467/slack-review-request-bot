package repository

import (
	"github.com/himura467/slack-review-request-bot/internal/domain/model"
	"net/http"
)

// SlackRepository defines the interface for Slack operations
type SlackRepository interface {
	// VerifyRequest validates the incoming request
	VerifyRequest(r *http.Request) error
	// ParseEvent parses the raw event data into a domain event
	ParseEvent(body []byte) (model.Event, error)
	// PostMessage posts a message to a Slack channel
	PostMessage(message *model.Message) error
}
