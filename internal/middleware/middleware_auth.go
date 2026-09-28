// Package middleware contains HTTP middleware shared by the API routes.
package middleware

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/ShortOwl/RSSAggregator/internal/auth"
	"github.com/ShortOwl/RSSAggregator/internal/database"
	"github.com/ShortOwl/RSSAggregator/internal/response"
)

// AuthedHandler is a handler that also receives the authenticated user.
type AuthedHandler func(http.ResponseWriter, *http.Request, database.User)

// WithAuth requires a Bearer JWT before calling handler.
func WithAuth(db *database.Queries, jwtSecret string, handler AuthedHandler) http.HandlerFunc {
	// returned handler ne yaad che db connection,jwtSecret ane authedhandler malse.
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Fields(r.Header.Get("Authorization"))

		if len(parts) != 2 {
			response.WithError(w, http.StatusUnauthorized, "Authentication credentials are required")
			return
		}

		if !strings.EqualFold(parts[0], "Bearer") {
			response.WithError(w, http.StatusUnauthorized, "Unsupported authentication scheme")
			return
		}
		userID, tokenError := auth.ValidateJWT(parts[1], jwtSecret)
		if tokenError != nil {
			response.WithError(w, http.StatusUnauthorized, "token is expired or invalid")
			return
		}
		user, err := db.GetUserByID(r.Context(), userID)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				response.WithError(w, http.StatusUnauthorized, "Invalid authentication credentials")
				return
			}
			response.WithError(w, http.StatusInternalServerError, "Couldn't authenticate user")
			return
		}
		handler(w, r, user)

	}

}
