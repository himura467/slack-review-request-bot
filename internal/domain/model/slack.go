package model

import "math/rand"

// OAuthToken represents a Slack OAuth token
type OAuthToken string

// SigningSecret represents a Slack signing secret
type SigningSecret string

// ReviewerIDs represents a list of Slack user IDs that can be assigned as reviewers
type ReviewerIDs []string

// GetRandomReviewer returns a random reviewer ID from the list
func (r ReviewerIDs) GetRandomReviewer() (string, bool) {
	if len(r) == 0 {
		return "", false
	}
	return r[rand.Intn(len(r))], true
}

// Message represents a Slack message
type Message struct {
	ChannelID string
	Text      string
}

func NewMessage(channelID, text string) *Message {
	return &Message{
		ChannelID: channelID,
		Text:      text,
	}
}

// Event represents a Slack event
type Event interface {
	Handle(handler EventHandler) *HTTPResponse
}

// EventHandler defines the interface for handling different types of events
type EventHandler interface {
	HandleAppMention(event *AppMentionEvent) *HTTPResponse
	HandleURLVerification(event *URLVerificationEvent) *HTTPResponse
}

// AppMentionEvent represents a Slack app mention event
type AppMentionEvent struct {
	ChannelID string
}

func NewAppMentionEvent(channelID string) *AppMentionEvent {
	return &AppMentionEvent{
		ChannelID: channelID,
	}
}

func (e *AppMentionEvent) Handle(handler EventHandler) *HTTPResponse {
	return handler.HandleAppMention(e)
}

// URLVerificationEvent represents a Slack URL verification event
type URLVerificationEvent struct {
	Challenge string
}

func NewURLVerificationEvent(challenge string) *URLVerificationEvent {
	return &URLVerificationEvent{
		Challenge: challenge,
	}
}

func (e *URLVerificationEvent) Handle(handler EventHandler) *HTTPResponse {
	return handler.HandleURLVerification(e)
}
