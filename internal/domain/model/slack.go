package model

import "math/rand"

// OAuthToken represents a Slack OAuth token
type OAuthToken string

// SigningSecret represents a Slack signing secret
type SigningSecret string

// MemberID represents a Slack member ID
type MemberID string

// ReviewerMap represents a mapping of display names to Slack member IDs for reviewers
type ReviewerMap map[string]MemberID

// Member contains both the display name and member ID of a Slack member
type Member struct {
	DisplayName string
	MemberID    MemberID
}

// GetRandomReviewer returns a random reviewer from the reviewer map.
// If filterMemberIDs is provided, it filters to only those member IDs.
// If excludeMemberIDs is provided, it excludes those member IDs from selection.
// If both are provided, it first filters then excludes.
func (r ReviewerMap) GetRandomReviewer(filterMemberIDs []MemberID, excludeMemberIDs []MemberID) (Member, bool) {
	if len(r) == 0 {
		return Member{}, false
	}

	// Create sets for efficient lookup
	var filterSet map[MemberID]bool
	if len(filterMemberIDs) > 0 {
		filterSet = make(map[MemberID]bool)
		for _, memberID := range filterMemberIDs {
			filterSet[memberID] = true
		}
	}
	var excludeSet map[MemberID]bool
	if len(excludeMemberIDs) > 0 {
		excludeSet = make(map[MemberID]bool)
		for _, memberID := range excludeMemberIDs {
			excludeSet[memberID] = true
		}
	}

	var candidates []Member
	for displayName, memberID := range r {
		// Apply filter if specified
		if filterSet != nil && !filterSet[memberID] {
			continue
		}
		// Apply exclusion if specified
		if excludeSet != nil && excludeSet[memberID] {
			continue
		}
		candidates = append(candidates, Member{
			DisplayName: displayName,
			MemberID:    memberID,
		})
	}
	if len(candidates) == 0 {
		return Member{}, false
	}
	// Select random candidate
	selectedReviewer := candidates[rand.Intn(len(candidates))]
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
	MemberID  MemberID
}

func NewInteractiveMessageEvent(channelID, actionID, value, messageTS, threadTS string, memberID MemberID) *InteractiveMessageEvent {
	return &InteractiveMessageEvent{
		ChannelID: channelID,
		ActionID:  actionID,
		Value:     value,
		MessageTS: messageTS,
		ThreadTS:  threadTS,
		MemberID:  memberID,
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
