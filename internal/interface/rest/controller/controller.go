package controller

import (
	"github.com/himura467/slack-review-request-bot/internal/usecase"
)

type Controller struct {
	slack usecase.SlackUsecase
}

func NewController(slack usecase.SlackUsecase) *Controller {
	return &Controller{
		slack: slack,
	}
}
