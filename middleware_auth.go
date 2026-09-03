package main

import (
	"fmt"
	"net/http"

	"github.com/ShortOwl/RSSAggregator/internal/auth"
	"github.com/ShortOwl/RSSAggregator/internal/database"
)

// a custom type for handlers that require authentication
type authHandler func(http.ResponseWriter, *http.Request, database.User)

// middleware that authenticates a request,gets the user and calls the next authed handler.

func (a *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)

		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("auth error:%v", err))
			return
		}

		user, err := a.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get user : %v", err))
			return
		}

		// AIya pochya etle user authenticate thai gayo, have custom handler run krai dav.
		// nice way of writing code.
		handler(w, r, user)
	}
}
