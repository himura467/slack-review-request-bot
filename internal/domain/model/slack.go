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

// Action represents a Slack message action
type Action struct {
	Name    string `json:"name"`
	Text    string `json:"text"`
	Type    string `json:"type"`
	Value   string `json:"value,omitempty"`
	Options []struct {
		Text  string `json:"text"`
		Value string `json:"value"`
	} `json:"options,omitempty"`
}

// Attachment represents a Slack message attachment
type Attachment struct {
	Text       string   `json:"text,omitempty"`
	CallbackID string   `json:"callback_id,omitempty"`
	Actions    []Action `json:"actions,omitempty"`
}

// Message represents a Slack message
type Message struct {
	ChannelID   string       `json:"channel"`
	Text        string       `json:"text,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// NewMessage creates a new Slack message
func NewMessage(channelID, text string) *Message {
	return &Message{
		ChannelID: channelID,
		Text:      text,
	}
}

// NewReviewerSelectionMessage creates a message with reviewer selection components using Slack Attachments
func NewReviewerSelectionMessage(channelID string, text string, reviewerMap ReviewerMap) *Message {
	// Create options for the select menu
	options := make([]struct {
		Text  string `json:"text"`
		Value string `json:"value"`
	}, 0, len(reviewerMap))
	for displayName, memberID := range reviewerMap {
		options = append(options, struct {
			Text  string `json:"text"`
			Value string `json:"value"`
		}{
			Text:  displayName,
			Value: memberID,
		})
	}

	return &Message{
		ChannelID: channelID,
		Attachments: []Attachment{
			{
				Text:       text,
				CallbackID: "reviewer_selection",
				Actions: []Action{
					{
						Name:  "random_reviewer",
						Text:  "Random",
						Type:  "button",
						Value: "",
					},
					{
						Name:    "select_reviewer",
						Text:    "レビュワーを選択",
						Type:    "select",
						Options: options,
					},
				},
			},
		},
	}
}

// Event represents a Slack event
type Event interface {
	Handle(handler EventHandler) *HTTPResponse
}

// EventHandler defines the interface for handling different types of events
type EventHandler interface {
	HandleAppMention(event *AppMentionEvent) *HTTPResponse
	HandleInteractiveMessage(event *InteractiveMessageEvent) *HTTPResponse
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

// InteractiveMessageEvent represents a Slack interactive message event
type InteractiveMessageEvent struct {
	ChannelID string
	ActionID  string
	Value     string
}

func NewInteractiveMessageEvent(channelID, actionID, value string) *InteractiveMessageEvent {
	return &InteractiveMessageEvent{
		ChannelID: channelID,
		ActionID:  actionID,
		Value:     value,
	}
}

func (e *InteractiveMessageEvent) Handle(handler EventHandler) *HTTPResponse {
	return handler.HandleInteractiveMessage(e)
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
