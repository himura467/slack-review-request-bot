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

// GetRandomOnlineReviewer returns a random online reviewer from the map
func (r ReviewerMap) GetRandomOnlineReviewer(onlineMembers []string) (ReviewerInfo, bool) {
	// Create a set of online member IDs for efficient lookup
	onlineMemberSet := make(map[string]bool)
	for _, memberID := range onlineMembers {
		onlineMemberSet[memberID] = true
	}
	// Get online reviewers
	var onlineReviewers []ReviewerInfo
	for displayName, memberID := range r {
		if onlineMemberSet[memberID] {
			onlineReviewers = append(onlineReviewers, ReviewerInfo{
				DisplayName: displayName,
				MemberID:    memberID,
			})
		}
	}
	if len(onlineReviewers) == 0 {
		return ReviewerInfo{}, false
	}
	// Select random online reviewer
	selectedReviewer := onlineReviewers[rand.Intn(len(onlineReviewers))]
	return selectedReviewer, true
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

// AttachmentField represents a field in a Slack message attachment
type AttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}

// Attachment represents a Slack message attachment
type Attachment struct {
	Text       string            `json:"text,omitempty"`
	Color      string            `json:"color,omitempty"`
	CallbackID string            `json:"callback_id,omitempty"`
	Actions    []Action          `json:"actions,omitempty"`
	Fields     []AttachmentField `json:"fields,omitempty"`
}

// Message represents a Slack message
type Message struct {
	ChannelID       string       `json:"channel"`
	Text            string       `json:"text,omitempty"`
	Attachments     []Attachment `json:"attachments,omitempty"`
	ReplaceOriginal bool         `json:"replace_original,omitempty"`
	ThreadTS        string       `json:"thread_ts,omitempty"`
}

func NewMessage(channelID, text string, attachments []Attachment, replaceOriginal bool, threadTS string) *Message {
	return &Message{
		ChannelID:       channelID,
		Text:            text,
		Attachments:     attachments,
		ReplaceOriginal: replaceOriginal,
		ThreadTS:        threadTS,
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
	ThreadTS  string
}

func NewAppMentionEvent(channelID string, threadTS string) *AppMentionEvent {
	return &AppMentionEvent{
		ChannelID: channelID,
		ThreadTS:  threadTS,
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
	MessageTS string
	ThreadTS  string
}

func NewInteractiveMessageEvent(channelID, actionID, value, messageTS, threadTS string) *InteractiveMessageEvent {
	return &InteractiveMessageEvent{
		ChannelID: channelID,
		ActionID:  actionID,
		Value:     value,
		MessageTS: messageTS,
		ThreadTS:  threadTS,
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
