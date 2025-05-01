package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/himura467/slack-review-request-bot/internal/usecase"
)

type app struct {
	slack usecase.SlackUsecase
}

func newApp(slack usecase.SlackUsecase) *app {
	return &app{
		slack: slack,
	}
}

func (a *app) Run() {
	r := chi.NewRouter()
	r.Post("/slack/events", a.slack.HandleEvent)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	slog.Info("starting server", "port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
