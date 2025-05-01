package usecase

import (
	"log/slog"
	"net/http"

	"github.com/himura467/slack-review-request-bot/internal/domain/model"
	"github.com/himura467/slack-review-request-bot/internal/domain/repository"
)

type SlackUsecase interface {
	HandleEvent(w http.ResponseWriter, r *http.Request)
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
func (u *SlackUsecaseImpl) HandleEvent(w http.ResponseWriter, r *http.Request) {
	// Verify the request and get body
	body, err := u.slackRepo.VerifyRequest(r)
	if err != nil {
		slog.Error("failed to verify request", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Parse the event
	event, err := u.slackRepo.ParseEvent(body)
	if err != nil {
		slog.Error("failed to parse event", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch e := event.(type) {
	case *model.CallbackEvent:
		// Only respond to non-threaded messages
		if e.IsThreadedMessage() {
			w.WriteHeader(http.StatusOK)
			return
		}
		// Get random reviewer from configured list
		reviewer, ok := u.reviewerIDs.GetRandomReviewer()
		if !ok {
			slog.Error("no reviewer IDs configured")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		messageText := "<@" + reviewer + "> このメッセージをレビューし、完了したら :white_check_mark: のリアクションをつけてください"
		message := model.NewMessage(e.GetChannelID(), messageText)
		// Post the message to Slack
		if err := u.slackRepo.PostMessage(message); err != nil {
			slog.Error("failed to post message", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	case *model.URLVerificationEvent:
		if _, err := w.Write([]byte(e.GetChallenge())); err != nil {
			slog.Error("failed to write response", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusOK)
		return
	}
}
