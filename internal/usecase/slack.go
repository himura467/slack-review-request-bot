package usecase

import (
	"log/slog"
	"net/http"

	"github.com/himura467/slack-review-request-bot/internal/domain/model"
	"github.com/himura467/slack-review-request-bot/internal/domain/repository"
)

type SlackUsecase interface {
	HandleEvent(r *model.HTTPRequest) *model.HTTPResponse
}

type SlackUsecaseImpl struct {
	slackRepo   repository.SlackRepository
	reviewerIDs model.ReviewerIDs
}

var _ SlackUsecase = (*SlackUsecaseImpl)(nil)

func NewSlackUsecase(slackRepo repository.SlackRepository, reviewerIDs model.ReviewerIDs) *SlackUsecaseImpl {
	return &SlackUsecaseImpl{
		slackRepo:   slackRepo,
		reviewerIDs: reviewerIDs,
	}
}

// HandleEvent processes incoming Slack events
func (u *SlackUsecaseImpl) HandleEvent(r *model.HTTPRequest) *model.HTTPResponse {
	// Verify the request
	if err := u.slackRepo.VerifyRequest(r); err != nil {
		slog.Error("failed to verify request", "error", err)
		return model.NewStatusResponse(http.StatusBadRequest)
	}
	// Parse the event
	event, err := u.slackRepo.ParseEvent(r.Body)
	if err != nil {
		slog.Error("failed to parse event", "error", err)
		return model.NewStatusResponse(http.StatusBadRequest)
	}
	switch e := event.(type) {
	case *model.CallbackEvent:
		// Only respond to non-threaded messages
		if e.IsThreadedMessage() {
			return model.NewStatusResponse(http.StatusOK)
		}
		// Get random reviewer from configured list
		reviewer, ok := u.reviewerIDs.GetRandomReviewer()
		if !ok {
			slog.Error("no reviewer IDs configured")
			return model.NewStatusResponse(http.StatusInternalServerError)
		}
		messageText := "<@" + reviewer + "> このメッセージをレビューし、完了したら :white_check_mark: のリアクションをつけてください"
		message := model.NewMessage(e.GetChannelID(), messageText)
		// Post the message to Slack
		if err := u.slackRepo.PostMessage(message); err != nil {
			slog.Error("failed to post message", "error", err)
			return model.NewStatusResponse(http.StatusInternalServerError)
		}
		return model.NewStatusResponse(http.StatusOK)
	case *model.URLVerificationEvent:
		return model.NewTextResponse(http.StatusOK, []byte(e.GetChallenge()))
	default:
		return model.NewStatusResponse(http.StatusOK)
	}
}
