package usecase

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/himura467/slack-review-request-bot/internal/domain/model"
)

// HandleAppMention handles app mention events
func (u *SlackUsecaseImpl) HandleAppMention(event *model.AppMentionEvent) *model.HTTPResponse {
	// Create options for the select menu
	options := make([]struct {
		Text  string `json:"text"`
		Value string `json:"value"`
	}, 0, len(u.reviewerMap))
	for displayName := range u.reviewerMap {
		options = append(options, struct {
			Text  string `json:"text"`
			Value string `json:"value"`
		}{
			Text:  displayName,
			Value: displayName,
		})
	}
	message := model.NewMessage(
		event.ChannelID,
		"レビュワーを選択してください",
		[]model.Attachment{
			{
				Text:       "ランダム指定もできるよ",
				CallbackID: "reviewer_selection",
				Actions: []model.Action{
					{
						Name:  "random_reviewer",
						Text:  "Random",
						Type:  "button",
						Value: "",
					},
					{
						Name:    "select_reviewer",
						Text:    "レビュワーを選択",
						Type:    "select",
						Options: options,
					},
				},
			},
		},
		false,
	)
	// Post the message to Slack
	if err := u.slackRepo.PostMessage(message); err != nil {
		slog.Error("failed to post message", "error", err)
		return model.NewStatusResponse(http.StatusInternalServerError)
	}
	return model.NewStatusResponse(http.StatusOK)
}

// HandleInteractiveMessage handles interactive message events
func (u *SlackUsecaseImpl) HandleInteractiveMessage(event *model.InteractiveMessageEvent) *model.HTTPResponse {
	var reviewerName string
	switch event.ActionID {
	case "random_reviewer":
		// Get random reviewer from configured map
		reviewer, ok := u.reviewerMap.GetRandomReviewer()
		if !ok {
			slog.Error("no reviewers configured")
			return model.NewStatusResponse(http.StatusInternalServerError)
		}
		reviewerName = reviewer.DisplayName
	case "select_reviewer":
		reviewerName = event.Value
	case "reassign_reviewer":
		// Get current reviewer name from Value field
		currentReviewerName := event.Value
		// Create a map of candidate reviewers excluding the current reviewer
		candidateReviewers := make(model.ReviewerMap)
		for name, id := range u.reviewerMap {
			if name != currentReviewerName {
				candidateReviewers[name] = id
			}
		}
		// Get random reviewer from the candidate reviewers
		reviewer, ok := candidateReviewers.GetRandomReviewer()
		if !ok {
			slog.Error("no other reviewers available")
			return model.NewStatusResponse(http.StatusInternalServerError)
		}
		reviewerName = reviewer.DisplayName
	default:
		slog.Error("unknown action ID", "action_id", event.ActionID)
		return model.NewStatusResponse(http.StatusBadRequest)
	}
	// Get reviewer ID for the mention
	reviewerID := u.reviewerMap[reviewerName]
	messageText := "<@" + reviewerID + "> このメッセージをレビューし、完了したら :white_check_mark: のリアクションをつけてください。\nメッセージ内のリンクは *シークレットウィンドウ* で開いて確認するようにしてください。"
	fields := []model.AttachmentField{
		{
			Title: "レビュワー",
			Value: reviewerName,
			Short: false,
		},
	}
	// Create Reassign button action
	actions := []model.Action{
		{
			Name:  "reassign_reviewer",
			Text:  "Reassign",
			Type:  "button",
			Value: reviewerName,
		},
	}
	message := model.NewMessage(
		event.ChannelID,
		messageText,
		[]model.Attachment{
			{
				Color:      "#F4631E",
				Fields:     fields,
				Actions:    actions,
				CallbackID: "reviewer_action",
			},
		},
		true,
	)
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
