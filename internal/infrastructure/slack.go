package infrastructure

import (
	"encoding/json"
	"log/slog"
	"net/url"

	"github.com/himura467/slack-review-request-bot/internal/domain/model"
	"github.com/himura467/slack-review-request-bot/internal/domain/repository"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type Client struct {
	api           *slack.Client
	signingSecret model.SigningSecret
}

var _ repository.SlackRepository = (*Client)(nil)

func NewClient(oauthToken model.OAuthToken, signingSecret model.SigningSecret) *Client {
	return &Client{
		api:           slack.New(string(oauthToken)),
		signingSecret: signingSecret,
	}
}

func (c *Client) VerifyRequest(r *model.HTTPRequest) error {
	sv, err := slack.NewSecretsVerifier(r.Headers, string(c.signingSecret))
	if err != nil {
		slog.Error("failed to create secrets verifier", "error", err)
		return err
	}
	if _, err = sv.Write(r.Body); err != nil {
		slog.Error("failed to write body to verifier", "error", err)
		return err
	}
	if err := sv.Ensure(); err != nil {
		slog.Error("failed to verify request", "error", err)
		return err
	}
	slog.Info("request verified successfully")
	return nil
}

func (c *Client) ParseEvent(body []byte) (model.Event, error) {
	// Parse regular Slack events
	eventsAPIEvent, err := slackevents.ParseEvent(body, slackevents.OptionNoVerifyToken())
	if err != nil {
		slog.Error("failed to parse event", "error", err)
		return nil, err
	}
	switch eventsAPIEvent.Type {
	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			// If the message is not in a thread, use its own timestamp as the thread timestamp
			threadTS := ev.ThreadTimeStamp
			if threadTS == "" {
				threadTS = ev.TimeStamp
			}
			return model.NewAppMentionEvent(ev.Channel, threadTS), nil
		default:
			slog.Info("unsupported inner event type", "type", ev)
			return nil, nil
		}
	case slackevents.URLVerification:
		var r *slackevents.ChallengeResponse
		if err := json.Unmarshal(body, &r); err != nil {
			slog.Error("failed to parse challenge", "error", err)
			return nil, err
		}
		return model.NewURLVerificationEvent(r.Challenge), nil
	default:
		slog.Info("unsupported event type", "type", eventsAPIEvent.Type)
		return nil, nil
	}
}

func (c *Client) ParseInteraction(body []byte) (model.Event, error) {
	payloadStr := string(body)
	if len(payloadStr) <= 8 || payloadStr[:8] != "payload=" {
		return nil, nil
	}
	// URL decode and remove "payload=" prefix
	decoded, err := url.QueryUnescape(payloadStr[8:])
	if err != nil {
		slog.Error("failed to unescape payload", "error", err)
		return nil, err
	}
	var interaction slack.InteractionCallback
	if err := json.Unmarshal([]byte(decoded), &interaction); err != nil {
		slog.Error("failed to parse interaction", "error", err)
		return nil, err
	}
	if len(interaction.ActionCallback.AttachmentActions) == 0 {
		return nil, nil
	}
	action := interaction.ActionCallback.AttachmentActions[0]
	var value string
	if action.Name == "random_reviewer" || action.Name == "urgent_reviewer" {
		value = "" // Empty value indicates random selection
	} else if action.Name == "select_reviewer" && len(action.SelectedOptions) > 0 {
		value = action.SelectedOptions[0].Value
	}
	// Get thread timestamp from the message
	threadTS := interaction.OriginalMessage.ThreadTimestamp
	if threadTS == "" {
		threadTS = interaction.OriginalMessage.Timestamp
	}
	return model.NewInteractiveMessageEvent(
		interaction.Channel.ID,
		action.Name,
		value,
		interaction.MessageTs,
		threadTS,
		model.MemberID(interaction.User.ID),
	), nil
}

func (c *Client) PostMessage(message *model.Message) error {
	var options []slack.MsgOption
	options = append(options, slack.MsgOptionText(message.Text, false))
	// When ThreadTS is set, ensure the message is posted in that thread
	if message.ThreadTS != "" {
		options = append(options, slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{
			ThreadTimestamp: message.ThreadTS,
		}))
	}

	if len(message.Attachments) > 0 {
		var attachments []slack.Attachment
		for _, a := range message.Attachments {
			var actions []slack.AttachmentAction
			for _, act := range a.Actions {
				action := slack.AttachmentAction{
					Name:  act.Name,
					Text:  act.Text,
					Type:  slack.ActionType(act.Type),
					Value: act.Value,
				}
				if len(act.Options) > 0 {
					actionOptions := make([]slack.AttachmentActionOption, len(act.Options))
					for i, opt := range act.Options {
						actionOptions[i] = slack.AttachmentActionOption{
							Text:  opt.Text,
							Value: opt.Value,
						}
					}
					action.Options = actionOptions
				}
				actions = append(actions, action)
			}
			attachment := slack.Attachment{
				Text:       a.Text,
				CallbackID: a.CallbackID,
				Actions:    actions,
				Color:      a.Color,
				Fields:     make([]slack.AttachmentField, len(a.Fields)),
			}
			// Convert Fields
			for i, f := range a.Fields {
				attachment.Fields[i] = slack.AttachmentField{
					Title: f.Title,
					Value: f.Value,
					Short: f.Short,
				}
			}
			attachments = append(attachments, attachment)
		}
		options = append(options, slack.MsgOptionAttachments(attachments...))
	}

	_, _, err := c.api.PostMessage(
		message.ChannelID,
		options...,
	)
	if err != nil {
		slog.Error("failed to post message", "error", err)
		return err
	}
	slog.Info("message posted successfully", "channel", message.ChannelID)
	return nil
}

func (c *Client) DeleteMessage(channelID, timestamp string) error {
	_, _, err := c.api.DeleteMessage(channelID, timestamp)
	if err != nil {
		slog.Error("failed to delete message", "error", err)
		return err
	}
	slog.Info("message deleted successfully", "channel", channelID)
	return nil
}

func (c *Client) FilterOnlineMemberIDs(memberIDs []model.MemberID) ([]model.MemberID, error) {
	var onlineMemberIDs []model.MemberID
	for _, memberID := range memberIDs {
		// Get user presence
		presence, err := c.api.GetUserPresence(string(memberID))
		if err != nil {
			slog.Warn("failed to get user presence", "user_id", memberID, "error", err)
			continue
		}
		// Check if user is active/online
		if presence.Presence == "active" {
			onlineMemberIDs = append(onlineMemberIDs, memberID)
		}
	}
	slog.Info("found online members", "input_count", len(memberIDs), "online_count", len(onlineMemberIDs))
	return onlineMemberIDs, nil
}
