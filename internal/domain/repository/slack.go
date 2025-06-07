package repository

import (
	"github.com/himura467/slack-review-request-bot/internal/domain/model"
)

// SlackRepository defines the interface for Slack operations
type SlackRepository interface {
	// VerifyRequest validates the incoming request
	VerifyRequest(r *model.HTTPRequest) error
	// ParseEvent parses the raw event data into a domain event
	ParseEvent(body []byte) (model.Event, error)
	// ParseInteraction parses the raw interaction data into a domain event
	ParseInteraction(body []byte) (model.Event, error)
	// PostMessage posts a message to a Slack channel
	PostMessage(message *model.Message) error
	// DeleteMessage deletes a message from a Slack channel
	DeleteMessage(channelID, timestamp string) error
	// FilterOnlineMemberIDs returns a list of online member IDs from the specified member IDs
	FilterOnlineMemberIDs(memberIDs []model.MemberID) ([]model.MemberID, error)
}
