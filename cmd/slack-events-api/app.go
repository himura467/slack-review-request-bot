package main

import (
	"log/slog"

	"github.com/himura467/slack-review-request-bot/internal/interface/rest"
)

type app struct {
	server *rest.Server
}

func newApp(server *rest.Server) *app {
	return &app{
		server: server,
	}
}

func (a *app) Run() {
	if err := a.server.Run(); err != nil {
		slog.Error("failed to run server", "error", err)
	}
}
