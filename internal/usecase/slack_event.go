package usecase

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/himura467/slack-review-request-bot/internal/domain/model"
)

// HandleAppMention handles app mention events
func (u *SlackUsecaseImpl) HandleAppMention(event *model.AppMentionEvent) *model.HTTPResponse {
	messageText := "レビュワーを選択してください"
	message := model.NewReviewerSelectionMessage(event.ChannelID, messageText, u.reviewerMap)
	// Post the message to Slack
	if err := u.slackRepo.PostMessage(message); err != nil {
		slog.Error("failed to post message", "error", err)
		return model.NewStatusResponse(http.StatusInternalServerError)
	}
	return model.NewStatusResponse(http.StatusOK)
}

// HandleInteractiveMessage handles interactive message events
func (u *SlackUsecaseImpl) HandleInteractiveMessage(event *model.InteractiveMessageEvent) *model.HTTPResponse {
	var reviewerID string
	switch event.ActionID {
	case "random_reviewer":
		// Get random reviewer from configured map
		reviewer, ok := u.reviewerMap.GetRandomReviewer()
		if !ok {
			slog.Error("no reviewers configured")
			return model.NewStatusResponse(http.StatusInternalServerError)
		}
		reviewerID = reviewer.MemberID
	case "select_reviewer":
		reviewerID = event.Value
	default:
		slog.Error("unknown action ID", "action_id", event.ActionID)
		return model.NewStatusResponse(http.StatusBadRequest)
	}
	messageText := "このメッセージをレビューし、完了したら :white_check_mark: のリアクションをつけてください。\nメッセージ内のリンクは *シークレットウィンドウ* で開いて確認するようにしてください。"
	fields := []model.AttachmentField{
		{
			Title: "レビュワー",
			Value: "<@" + reviewerID + ">",
			Short: false,
		},
	}
	message := model.NewUpdateMessage(event.ChannelID, messageText, fields)
	// Encode response as JSON
	responseJSON, err := json.Marshal(message)
	if err != nil {
		slog.Error("failed to marshal response", "error", err)
		return model.NewStatusResponse(http.StatusInternalServerError)
	}
	return model.NewJSONResponse(http.StatusOK, responseJSON)
}

// HandleURLVerification handles URL verification events
func (u *SlackUsecaseImpl) HandleURLVerification(event *model.URLVerificationEvent) *model.HTTPResponse {
	return model.NewTextResponse(http.StatusOK, []byte(event.Challenge))
}
