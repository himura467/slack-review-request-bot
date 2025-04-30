package infrastructure

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/himura467/slack-review-request-bot/internal/domain/model"
	"github.com/himura467/slack-review-request-bot/internal/domain/repository"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type Client struct {
	api           *slack.Client
	signingSecret string
}

var _ repository.SlackRepository = (*Client)(nil)

func NewClient(oauthToken, signingSecret string) *Client {
	return &Client{
		api:           slack.New(oauthToken),
		signingSecret: signingSecret,
	}
}

func (c *Client) VerifyRequest(r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("failed to read request body", "error", err)
		return err
	}
	sv, err := slack.NewSecretsVerifier(r.Header, c.signingSecret)
	if err != nil {
		slog.Error("failed to create secrets verifier", "error", err)
		return err
	}
	if _, err = sv.Write(body); err != nil {
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
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
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
		return model.NewURLVerificationEvent(eventsAPIEvent.Type, r.Challenge), nil
	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			if ev.BotID != "" {
				slog.Info("ignoring bot message", "bot_id", ev.BotID)
				return nil, nil
			}
			return model.NewCallbackEvent(eventsAPIEvent.Type, ev.Channel, ev.ThreadTimeStamp), nil
		default:
			slog.Error("unsupported inner event type", "type", ev)
			return nil, nil
		}
	default:
		slog.Error("unsupported event type", "type", eventsAPIEvent.Type)
		return nil, err
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
