package handler

import (
	"net/http"

	"github.com/ShortOwl/RSSAggregator/internal/response"
)

// HandleReadiness is a health check endpoint.
// Load balancers and monitoring tools call this to verify the server is alive.
func (h *Handler) HandleReadiness(w http.ResponseWriter, r *http.Request) {
	response.WithJSON(w, 200, struct{}{})
}

// HandleError is a test endpoint that always returns an error.
// Useful during development to verify error handling works.
func (h *Handler) HandleError(w http.ResponseWriter, r *http.Request) {
	response.WithError(w, 400, "Something went wrong :(")
}
