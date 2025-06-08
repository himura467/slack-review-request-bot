package usecase

import (
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
				Text:       "ランダムに指定したい場合は「Random」を、急ぎの場合は「Urgent」を選択してください",
				CallbackID: "reviewer_selection",
				Actions: []model.Action{
					{
						Name:  "random_reviewer",
						Text:  "Random",
						Type:  "button",
						Value: "",
					},
					{
						Name:  "urgent_reviewer",
						Text:  "Urgent",
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
		event.ThreadTS,
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
	var messageText string
	switch event.ActionID {
	case "random_reviewer":
		// Get random reviewer from configured map
		reviewer, ok := u.reviewerMap.GetRandomReviewerFrom(nil)
		if !ok {
			slog.Error("no reviewers configured")
			return model.NewStatusResponse(http.StatusInternalServerError)
		}
		reviewerName = reviewer.DisplayName
		reviewerID := u.reviewerMap[reviewerName]
		messageText = "<@" + string(reviewerID) + ">\n【ランダム】\nこのメッセージをレビューし、完了したら :white_check_mark: のリアクションをつけてください。\nメッセージ内のリンクは *シークレットウィンドウ* で開いて確認するようにしてください。"
	case "urgent_reviewer":
		// Process urgent reviewer selection asynchronously to avoid timeout
		go u.processUrgentReviewer(event)
		// Return immediately to avoid Slack timeout
		return model.NewStatusResponse(http.StatusOK)
	case "select_reviewer":
		reviewerName = event.Value
		reviewerID := u.reviewerMap[reviewerName]
		messageText = "<@" + string(reviewerID) + ">\n【選択】\nこのメッセージをレビューし、完了したら :white_check_mark: のリアクションをつけてください。\nメッセージ内のリンクは *シークレットウィンドウ* で開いて確認するようにしてください。"
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
		reviewer, ok := candidateReviewers.GetRandomReviewerFrom(nil)
		if !ok {
			slog.Error("no other reviewers available")
			return model.NewStatusResponse(http.StatusInternalServerError)
		}
		reviewerName = reviewer.DisplayName
		reviewerID := u.reviewerMap[reviewerName]
		messageText = "<@" + string(reviewerID) + ">\n【ランダム】\nこのメッセージをレビューし、完了したら :white_check_mark: のリアクションをつけてください。\nメッセージ内のリンクは *シークレットウィンドウ* で開いて確認するようにしてください。"
	default:
		slog.Error("unknown action ID", "action_id", event.ActionID)
		return model.NewStatusResponse(http.StatusBadRequest)
	}
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
	// Delete the original message
	if err := u.slackRepo.DeleteMessage(event.ChannelID, event.MessageTS); err != nil {
		slog.Error("failed to delete message", "error", err)
		return model.NewStatusResponse(http.StatusInternalServerError)
	}
	// Create and post new message in the thread
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
		false,
		event.ThreadTS,
	)
	// Post the new message
	if err := u.slackRepo.PostMessage(message); err != nil {
		slog.Error("failed to post message", "error", err)
		return model.NewStatusResponse(http.StatusInternalServerError)
	}
	return model.NewStatusResponse(http.StatusOK)
}

// processUrgentReviewer handles urgent reviewer selection asynchronously
func (u *SlackUsecaseImpl) processUrgentReviewer(event *model.InteractiveMessageEvent) {
	// Get all reviewer member IDs from the map
	var allReviewerIDs []model.MemberID
	for _, memberID := range u.reviewerMap {
		allReviewerIDs = append(allReviewerIDs, memberID)
	}
	// Filter to get online member IDs from all reviewers
	onlineMemberIDs, err := u.slackRepo.FilterOnlineMemberIDs(allReviewerIDs)
	if err != nil {
		slog.Error("failed to filter online member IDs", "error", err)
		return
	}
	// Get random online reviewer from configured map
	reviewer, ok := u.reviewerMap.GetRandomReviewerFrom(onlineMemberIDs)
	if !ok {
		slog.Error("no reviewers configured")
		return
	}
	reviewerName := reviewer.DisplayName
	reviewerID := u.reviewerMap[reviewerName]
	messageText := "<@" + string(reviewerID) + ">\n【急ぎ】\nこのメッセージをレビューし、完了したら :white_check_mark: のリアクションをつけてください。\nメッセージ内のリンクは *シークレットウィンドウ* で開いて確認するようにしてください。"

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
	// Delete the original message
	if err := u.slackRepo.DeleteMessage(event.ChannelID, event.MessageTS); err != nil {
		slog.Error("failed to delete message", "error", err)
		return
	}
	// Create and post new message in the thread
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
		false,
		event.ThreadTS,
	)
	// Post the new message
	if err := u.slackRepo.PostMessage(message); err != nil {
		slog.Error("failed to post message", "error", err)
		return
	}
}

// HandleURLVerification handles URL verification events
func (u *SlackUsecaseImpl) HandleURLVerification(event *model.URLVerificationEvent) *model.HTTPResponse {
	return model.NewTextResponse(http.StatusOK, []byte(event.Challenge))
}
