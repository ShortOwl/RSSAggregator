package handler

import (
	"net/http"

	"github.com/ShortOwl/RSSAggregator/internal/database"
	"github.com/ShortOwl/RSSAggregator/internal/model"
	"github.com/ShortOwl/RSSAggregator/internal/response"
)

// HandleGetUser returns the currently authenticated user's info.
func (h *Handler) HandleGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	response.WithJSON(w, 200, model.DatabaseUserToUser(user))
}
