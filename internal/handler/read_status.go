package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/ShortOwl/RSSAggregator/internal/database"
	"github.com/ShortOwl/RSSAggregator/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func postIDFromURL(r *http.Request) (uuid.UUID, error) {
	return uuid.Parse(chi.URLParam(r, "postID"))
}

// HandleMarkPostAsRead marks a post as read for the authenticated user.
func (h *Handler) HandleMarkPostAsRead(w http.ResponseWriter, r *http.Request, user database.User) {
	postID, err := postIDFromURL(r)
	if err != nil {
		response.WithError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	_, err = h.DB.MarkPostAsRead(r.Context(), database.MarkPostAsReadParams{
		UserID: user.ID,
		PostID: postID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		response.WithError(w, http.StatusNotFound, "Post not found")
		return
	}
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Couldn't mark post as read")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleMarkPostAsUnread marks a post as unread for the authenticated user.
func (h *Handler) HandleMarkPostAsUnread(w http.ResponseWriter, r *http.Request, user database.User) {
	postID, err := postIDFromURL(r)
	if err != nil {
		response.WithError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	err = h.DB.MarkPostAsUnread(r.Context(), database.MarkPostAsUnreadParams{
		UserID: user.ID,
		PostID: postID,
	})
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Couldn't mark post as unread")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
