//go:build wireinject
// +build wireinject

package infrastructure

import (
	"github.com/google/wire"
	"github.com/himura467/slack-review-request-bot/internal/domain/repository"
)

var Set = wire.NewSet(
	NewClient,
	wire.Bind(new(repository.SlackRepository), new(*Client)),
)
