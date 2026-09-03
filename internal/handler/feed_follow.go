package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ShortOwl/RSSAggregator/internal/database"
	"github.com/ShortOwl/RSSAggregator/internal/model"
	"github.com/ShortOwl/RSSAggregator/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// HandleCreateFeedFollow lets an authenticated user follow a feed.
func (h *Handler) HandleCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		response.WithError(w, http.StatusBadRequest, "Couldn't decode parameters")
		return
	}

	feedFollow, err := h.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}

	response.WithJSON(w, http.StatusCreated, model.DatabaseFeedFollowToFeedFollow(feedFollow))
}

// HandleGetFeedFollows returns all feeds the authenticated user follows.
func (h *Handler) HandleGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := h.DB.GetFeedFollowsForUser(r.Context(), user.ID)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Couldn't get feed follows")
		return
	}

	response.WithJSON(w, http.StatusOK, model.DatabaseFeedFollowsToFeedFollows(feedFollows))
}

// HandleDeleteFeedFollow lets an authenticated user unfollow a feed.
func (h *Handler) HandleDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		response.WithError(w, http.StatusBadRequest, "Invalid feed follow ID")
		return
	}

	err = h.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Couldn't delete feed follow")
		return
	}

	response.WithJSON(w, http.StatusNoContent, struct{}{})
}
