//go:build wireinject
// +build wireinject

package controller

import (
	"github.com/google/wire"
	"github.com/himura467/slack-review-request-bot/internal/usecase"
)

var Set = wire.NewSet(
	usecase.Set,
	NewController,
)
