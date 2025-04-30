package usecase

import (
	"net/http"

	"github.com/himura467/slack-review-request-bot/internal/domain/model"
	"github.com/himura467/slack-review-request-bot/internal/domain/repository"
)

type SlackUsecase interface {
	HandleEvent(w http.ResponseWriter, r *http.Request) error
}

type SlackUsecaseImpl struct {
	slackRepo repository.SlackRepository
}

var _ SlackUsecase = (*SlackUsecaseImpl)(nil)

func NewSlackUsecase(slackRepo repository.SlackRepository) *SlackUsecaseImpl {
	return &SlackUsecaseImpl{
		slackRepo: slackRepo,
	}
}

// HandleEvent processes incoming Slack events
func (u *SlackUsecaseImpl) HandleEvent(w http.ResponseWriter, r *http.Request) error {
	// Verify the request and get body
	body, err := u.slackRepo.VerifyRequest(r)
	if err != nil {
		return err
	}
	// Parse the event
	event, err := u.slackRepo.ParseEvent(body)
	if err != nil {
		return err
	}
	switch e := event.(type) {
	case *model.CallbackEvent:
		// Only respond to non-threaded messages
		if e.IsThreadedMessage() {
			return nil
		}
		// Create a new message with "Hello World"
		message := model.NewMessage(e.GetChannelID(), "Hello World")
		// Post the message to Slack
		return u.slackRepo.PostMessage(message)
	case *model.URLVerificationEvent:
		// URL verification events don't require a response here
		return nil
	default:
		return nil
	}
}
