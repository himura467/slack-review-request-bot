package usecase

import (
	"log/slog"
	"net/http"

	"github.com/himura467/slack-review-request-bot/internal/domain/model"
)

// HandleAppMention handles app mention events
func (u *SlackUsecaseImpl) HandleAppMention(event *model.AppMentionEvent) *model.HTTPResponse {
	// Get random reviewer from configured map
	reviewer, ok := u.reviewerMap.GetRandomReviewer()
	if !ok {
		slog.Error("no reviewers configured")
		return model.NewStatusResponse(http.StatusInternalServerError)
	}
	messageText := "<@" + reviewer.MemberID + "> このメッセージをレビューし、完了したら :white_check_mark: のリアクションをつけてください"
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
