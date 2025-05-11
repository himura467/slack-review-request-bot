//go:build wireinject
// +build wireinject

package rest

import (
	"github.com/google/wire"
	"github.com/himura467/slack-review-request-bot/internal/interface/rest/controller"
)

var Set = wire.NewSet(
	controller.Set,
	NewServer,
)
