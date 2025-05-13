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
			return model.NewAppMentionEvent(ev.Channel), nil
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
	if interaction.Type != slack.InteractionTypeBlockActions {
		return nil, nil
	}
	if len(interaction.ActionCallback.BlockActions) == 0 {
		return nil, nil
	}
	action := interaction.ActionCallback.BlockActions[0]
	var value string
	if action.ActionID == "random_reviewer" {
		value = "" // Empty value indicates random selection
	} else {
		value = action.SelectedOption.Value
	}
	return model.NewInteractiveMessageEvent(interaction.Channel.ID, action.ActionID, value), nil
}

func (c *Client) PostMessage(message *model.Message) error {
	var options []slack.MsgOption
	options = append(options, slack.MsgOptionText(message.Text, false))

	if len(message.Blocks) > 0 {
		var blocks []slack.Block
		for _, b := range message.Blocks {
			switch b.Type {
			case "section":
				blocks = append(blocks, slack.NewSectionBlock(
					&slack.TextBlockObject{
						Type: b.Text.Type,
						Text: b.Text.Text,
					},
					nil,
					nil,
				))
			case "actions":
				var elements []slack.BlockElement
				for _, e := range b.Elements {
					switch e.Type {
					case "button":
						elements = append(elements, slack.NewButtonBlockElement(
							e.ActionID,
							e.ActionID,
							&slack.TextBlockObject{
								Type: e.Text.Type,
								Text: e.Text.Text,
							},
						))
					case "static_select":
						var options []*slack.OptionBlockObject
						for _, o := range e.Options {
							options = append(options, slack.NewOptionBlockObject(
								o.Value,
								&slack.TextBlockObject{
									Type: o.Text.Type,
									Text: o.Text.Text,
								},
								nil,
							))
						}
						elements = append(elements, slack.NewOptionsSelectBlockElement(
							slack.OptTypeStatic,
							&slack.TextBlockObject{
								Type: e.Placeholder.Type,
								Text: e.Placeholder.Text,
							},
							e.ActionID,
							options...,
						))
					}
				}
				blocks = append(blocks, slack.NewActionBlock(b.BlockID, elements...))
			}
		}
		options = append(options, slack.MsgOptionBlocks(blocks...))
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
