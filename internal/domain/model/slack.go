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

// Text represents text in Slack Block Kit
type Text struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Option represents an option in a Slack Block Kit select menu
type Option struct {
	Text  Text   `json:"text"`
	Value string `json:"value"`
}

// Element represents an interactive element in Slack Block Kit
type Element struct {
	Type        string   `json:"type"`
	ActionID    string   `json:"action_id,omitempty"`
	Text        *Text    `json:"text,omitempty"`
	Options     []Option `json:"options,omitempty"`
	Placeholder *Text    `json:"placeholder,omitempty"`
}

// Block represents a Slack Block Kit block
type Block struct {
	Type     string    `json:"type"`
	Text     *Text     `json:"text,omitempty"`
	Elements []Element `json:"elements,omitempty"`
	BlockID  string    `json:"block_id,omitempty"`
}

// Message represents a Slack message
type Message struct {
	ChannelID string  `json:"channel"`
	Text      string  `json:"text"`
	Blocks    []Block `json:"blocks,omitempty"`
}

// NewMessage creates a new Slack message
func NewMessage(channelID, text string) *Message {
	return &Message{
		ChannelID: channelID,
		Text:      text,
	}
}

// NewReviewerSelectionMessage creates a message with reviewer selection components using Slack Block Kit
func NewReviewerSelectionMessage(channelID string, text string, reviewerMap ReviewerMap) *Message {
	// Create options for the select menu
	options := make([]Option, 0, len(reviewerMap))
	for displayName, memberID := range reviewerMap {
		options = append(options, Option{
			Text: Text{
				Type: "plain_text",
				Text: displayName,
			},
			Value: memberID,
		})
	}

	return &Message{
		ChannelID: channelID,
		Text:      text,
		Blocks: []Block{
			{
				Type: "section",
				Text: &Text{
					Type: "mrkdwn",
					Text: text,
				},
			},
			{
				Type:    "actions",
				BlockID: "reviewer_selection",
				Elements: []Element{
					{
						Type:     "button",
						ActionID: "random_reviewer",
						Text: &Text{
							Type: "plain_text",
							Text: "Random",
						},
					},
					{
						Type:     "static_select",
						ActionID: "select_reviewer",
						Placeholder: &Text{
							Type: "plain_text",
							Text: "レビュワーを選択",
						},
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
