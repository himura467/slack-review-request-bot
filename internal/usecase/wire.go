//go:build wireinject
// +build wireinject

package usecase

import (
	"github.com/google/wire"
	"github.com/himura467/slack-review-request-bot/internal/infrastructure"
)

var Set = wire.NewSet(
	infrastructure.Set,
	NewSlackUsecase,
	wire.Bind(new(SlackUsecase), new(*SlackUsecaseImpl)),
)
