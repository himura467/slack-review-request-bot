package model

// Message represents a Slack message
type Message struct {
	ChannelID string
	Text      string
}

// NewMessage creates a new Message instance
func NewMessage(channelID, text string) *Message {
	return &Message{
		ChannelID: channelID,
		Text:      text,
	}
}

// Event represents a Slack event interface
type Event interface {
	GetType() string
}

// CallbackEvent represents a Slack callback event
type CallbackEvent struct {
	eventType string
	channelID string
	threadTS  string
}

func NewCallbackEvent(eventType, channelID, threadTS string) *CallbackEvent {
	return &CallbackEvent{
		eventType: eventType,
		channelID: channelID,
		threadTS:  threadTS,
	}
}

func (e *CallbackEvent) GetType() string {
	return e.eventType
}

// GetChannelID returns the channel ID of the event
func (e *CallbackEvent) GetChannelID() string {
	return e.channelID
}

// IsThreadedMessage checks if the message is part of a thread
func (e *CallbackEvent) IsThreadedMessage() bool {
	return e.threadTS != ""
}

// IsFromChannel checks if the message is from the specified channel
func (e *CallbackEvent) IsFromChannel(channelID string) bool {
	return e.channelID == channelID
}

// URLVerificationEvent represents a Slack URL verification event
type URLVerificationEvent struct {
	eventType string
	challenge string
}

func NewURLVerificationEvent(eventType, challenge string) *URLVerificationEvent {
	return &URLVerificationEvent{
		eventType: eventType,
		challenge: challenge,
	}
}

func (e *URLVerificationEvent) GetType() string {
	return e.eventType
}

func (e *URLVerificationEvent) GetChallenge() string {
	return e.challenge
}
