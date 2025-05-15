package rest

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/himura467/slack-review-request-bot/internal/interface/rest/controller"
)

type Server struct {
	router     *chi.Mux
	controller *controller.Controller
}

func NewServer(controller *controller.Controller) *Server {
	return &Server{
		router:     chi.NewRouter(),
		controller: controller,
	}
}

func (s *Server) Run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	s.router.Post("/slack/events", s.controller.HandleEvent)
	s.router.Post("/slack/interactions", s.controller.HandleInteraction)

	return http.ListenAndServe(":"+port, s.router)
}
