package model

import "math/rand"

// OAuthToken represents a Slack OAuth token
type OAuthToken string

// SigningSecret represents a Slack signing secret
type SigningSecret string

// ReviewerMap represents a mapping of display names to Slack member IDs for reviewers
type ReviewerMap map[string]string

// ReviewerInfo contains both the display name and member ID of a reviewer
type ReviewerInfo struct {
	DisplayName string
	MemberID    string
}

// GetRandomReviewer returns a random reviewer with both display name and member ID from the map
func (r ReviewerMap) GetRandomReviewer() (ReviewerInfo, bool) {
	if len(r) == 0 {
		return ReviewerInfo{}, false
	}
	// Get all display names as slice
	displayNames := make([]string, 0, len(r))
	for name := range r {
		displayNames = append(displayNames, name)
	}
	// Select random display name
	selectedName := displayNames[rand.Intn(len(displayNames))]
	return ReviewerInfo{
		DisplayName: selectedName,
		MemberID:    r[selectedName],
	}, true
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
