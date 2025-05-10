package controller

import "net/http"

func (c *Controller) HandleEvent(w http.ResponseWriter, r *http.Request) {
	c.slack.HandleEvent(w, r)
}
