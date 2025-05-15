package usecase

import (
	"log/slog"
	"net/http"

	"github.com/himura467/slack-review-request-bot/internal/domain/model"
	"github.com/himura467/slack-review-request-bot/internal/domain/repository"
)

type SlackUsecase interface {
	HandleEvent(r *model.HTTPRequest) *model.HTTPResponse
	HandleInteraction(r *model.HTTPRequest) *model.HTTPResponse
}

type SlackUsecaseImpl struct {
	slackRepo   repository.SlackRepository
	reviewerMap model.ReviewerMap
}

var _ SlackUsecase = (*SlackUsecaseImpl)(nil)
var _ model.EventHandler = (*SlackUsecaseImpl)(nil)

func NewSlackUsecase(slackRepo repository.SlackRepository, reviewerMap model.ReviewerMap) *SlackUsecaseImpl {
	return &SlackUsecaseImpl{
		slackRepo:   slackRepo,
		reviewerMap: reviewerMap,
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
	if event == nil {
		return model.NewStatusResponse(http.StatusOK)
	}
	return event.Handle(u)
}

// HandleInteraction processes incoming Slack interactions
func (u *SlackUsecaseImpl) HandleInteraction(r *model.HTTPRequest) *model.HTTPResponse {
	// Verify the request
	if err := u.slackRepo.VerifyRequest(r); err != nil {
		slog.Error("failed to verify request", "error", err)
		return model.NewStatusResponse(http.StatusBadRequest)
	}
	// Parse the interaction
	event, err := u.slackRepo.ParseInteraction(r.Body)
	if err != nil {
		slog.Error("failed to parse interaction", "error", err)
		return model.NewStatusResponse(http.StatusBadRequest)
	}
	if event == nil {
		return model.NewStatusResponse(http.StatusOK)
	}
	return event.Handle(u)
}
