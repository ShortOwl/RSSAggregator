package handler

import (
	"log"
	"net/http"

	"github.com/ShortOwl/RSSAggregator/internal/database"
	"github.com/ShortOwl/RSSAggregator/internal/model"
	"github.com/ShortOwl/RSSAggregator/internal/response"
)

// HandleGetPostsForUser returns posts from feeds the authenticated user follows.
func (h *Handler) HandleGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := h.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		log.Println("Error fetching posts for user:", user.ID, "err:", err)
		response.WithError(w, 500, "Couldn't fetch posts")
		return
	}

	response.WithJSON(w, 200, model.DatabasePostsToPosts(posts))
}
