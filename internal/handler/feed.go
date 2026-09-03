package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ShortOwl/RSSAggregator/internal/database"
	"github.com/ShortOwl/RSSAggregator/internal/model"
	"github.com/ShortOwl/RSSAggregator/internal/response"
	"github.com/google/uuid"
)

// HandleCreateFeed adds a new RSS feed to the system.
// Requires authentication — the user is passed in by WithAuth.
func (h *Handler) HandleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameter struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)
	if err != nil {
		response.WithError(w, 400, "Error parsing JSON body")
		return
	}

	feed, err := h.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		response.WithError(w, 500, "Couldn't create feed")
		return
	}

	response.WithJSON(w, 201, model.DatabaseFeedToFeed(feed))
}

// HandleGetFeeds returns all feeds in the system.
// This is a public endpoint — no authentication required.
func (h *Handler) HandleGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := h.DB.GetFeeds(r.Context())
	if err != nil {
		response.WithError(w, 500, "Error fetching feeds")
		return
	}

	response.WithJSON(w, 200, model.DatabaseFeedsToFeeds(feeds))
}
