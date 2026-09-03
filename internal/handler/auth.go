package handler

import (
	"fmt"
	"net/http"

	"github.com/ShortOwl/RSSAggregator/internal/auth"
	"github.com/ShortOwl/RSSAggregator/internal/database"
	"github.com/ShortOwl/RSSAggregator/internal/response"
)

// AuthedHandler is a handler that also receives the authenticated user.
type AuthedHandler func(http.ResponseWriter, *http.Request, database.User)

// WithAuth wraps an authenticated handler.
// It extracts the API key, looks up the user, and passes it to your handler.
func (h *Handler) WithAuth(handler AuthedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			response.WithError(w, 403, fmt.Sprintf("auth error: %v", err))
			return
		}

		user, err := h.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			response.WithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}

		// User is authenticated — call the actual handler with the user
		handler(w, r, user)
	}
}
