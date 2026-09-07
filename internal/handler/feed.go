package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
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

// HandleGetFeeds returns a page of feeds in the system.
// This is a public endpoint — no authentication required.
func (h *Handler) HandleGetFeeds(w http.ResponseWriter, r *http.Request) {
	limit, err := parseLimit(r)
	if err != nil {
		response.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	cursor, err := parseCursor(r)

	if err != nil {
		response.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	search := r.URL.Query().Get("search")

	params := database.GetFeedsParams{
		Search: strings.TrimSpace(search),
		Limit:  limit + 1,
	}

	if cursor != nil {
		params.CursorCreatedAt = sql.NullTime{Time: cursor.Timestamp, Valid: true}
		params.CursorID = uuid.NullUUID{UUID: cursor.ID, Valid: true}
	}

	feeds, err := h.DB.GetFeeds(r.Context(), params)
	if err != nil {
		response.WithError(w, 500, "Error fetching feeds")
		return
	}

	nextCursor := ""
	if len(feeds) > int(limit) {
		feeds = feeds[:limit]
		last := feeds[len(feeds)-1]
		nextCursor = encodeCursor(last.CreatedAt, last.ID)
	}

	response.WithJSON(w, http.StatusOK, struct {
		Data       []model.Feed `json:"data"`
		NextCursor string       `json:"next_cursor"`
	}{
		Data:       model.DatabaseFeedsToFeeds(feeds),
		NextCursor: nextCursor,
	})
}
