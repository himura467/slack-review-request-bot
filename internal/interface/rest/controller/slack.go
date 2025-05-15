package controller

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/himura467/slack-review-request-bot/internal/domain/model"
)

func (c *Controller) HandleEvent(w http.ResponseWriter, r *http.Request) {
	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("failed to read request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Create HTTPRequest
	request := model.NewHTTPRequest(body, r.Header)
	// Process the event through usecase
	response := c.slack.HandleEvent(request)
	// Set response content type if specified
	if response.ContentType != "" {
		w.Header().Set("Content-Type", response.ContentType)
	}
	// Set status code
	w.WriteHeader(response.StatusCode)
	// Write response body if present
	if len(response.Body) > 0 {
		if _, err := w.Write(response.Body); err != nil {
			slog.Error("failed to write response", "error", err)
			return
		}
	}
}

func (c *Controller) HandleInteraction(w http.ResponseWriter, r *http.Request) {
	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("failed to read request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Create HTTPRequest
	request := model.NewHTTPRequest(body, r.Header)
	// Process the interaction through usecase
	response := c.slack.HandleInteraction(request)
	// Set response content type if specified
	if response.ContentType != "" {
		w.Header().Set("Content-Type", response.ContentType)
	}
	// Set status code
	w.WriteHeader(response.StatusCode)
	// Write response body if present
	if len(response.Body) > 0 {
		if _, err := w.Write(response.Body); err != nil {
			slog.Error("failed to write response", "error", err)
			return
		}
	}
}
