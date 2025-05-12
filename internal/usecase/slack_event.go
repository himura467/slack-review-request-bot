package usecase

import (
	"log/slog"
	"net/http"

	"github.com/himura467/slack-review-request-bot/internal/domain/model"
)

// HandleCallback handles callback events
func (u *SlackUsecaseImpl) HandleCallback(event *model.CallbackEvent) *model.HTTPResponse {
	// Only respond to non-threaded messages
	if event.IsThreadedMessage() {
		return model.NewStatusResponse(http.StatusOK)
	}
	// Get random reviewer from configured list
	reviewer, ok := u.reviewerIDs.GetRandomReviewer()
	if !ok {
		slog.Error("no reviewer IDs configured")
		return model.NewStatusResponse(http.StatusInternalServerError)
	}
	messageText := "<@" + reviewer + "> このメッセージをレビューし、完了したら :white_check_mark: のリアクションをつけてください"
	message := model.NewMessage(event.ChannelID, messageText)
	// Post the message to Slack
	if err := u.slackRepo.PostMessage(message); err != nil {
		slog.Error("failed to post message", "error", err)
		return model.NewStatusResponse(http.StatusInternalServerError)
	}
	return model.NewStatusResponse(http.StatusOK)
}

// HandleURLVerification handles URL verification events
func (u *SlackUsecaseImpl) HandleURLVerification(event *model.URLVerificationEvent) *model.HTTPResponse {
	return model.NewTextResponse(http.StatusOK, []byte(event.Challenge))
}
