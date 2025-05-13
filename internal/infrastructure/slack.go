package infrastructure

import (
	"encoding/json"
	"log/slog"

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
	eventsAPIEvent, err := slackevents.ParseEvent(body, slackevents.OptionNoVerifyToken())
	if err != nil {
		slog.Error("failed to parse event", "error", err)
		return nil, err
	}
	switch eventsAPIEvent.Type {
	case slackevents.URLVerification:
		var r *slackevents.ChallengeResponse
		if err := json.Unmarshal(body, &r); err != nil {
			slog.Error("failed to parse challenge", "error", err)
			return nil, err
		}
		return model.NewURLVerificationEvent(r.Challenge), nil
	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			return model.NewAppMentionEvent(ev.Channel), nil
		default:
			slog.Info("unsupported inner event type", "type", ev)
			return nil, nil
		}
	default:
		slog.Info("unsupported event type", "type", eventsAPIEvent.Type)
		return nil, nil
	}
}

func (c *Client) PostMessage(message *model.Message) error {
	_, _, err := c.api.PostMessage(
		message.ChannelID,
		slack.MsgOptionText(message.Text, false),
	)
	if err != nil {
		slog.Error("failed to post message", "error", err)
		return err
	}
	slog.Info("message posted successfully", "channel", message.ChannelID)
	return nil
}
