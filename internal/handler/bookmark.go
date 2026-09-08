// akhu jate lakhva nu che.

package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/ShortOwl/RSSAggregator/internal/database"
	"github.com/ShortOwl/RSSAggregator/internal/model"
	"github.com/ShortOwl/RSSAggregator/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// all handlers are authenticated.
func (h *Handler) HandleBookmarkPost(w http.ResponseWriter, r *http.Request, user database.User) {
	postID, err := uuid.Parse(chi.URLParam(r, "postID"))

	if err != nil {

		response.WithError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	params := database.CreateBookmarkParams{
		UserID: user.ID,
		PostID: postID,
	}

	_, err = h.DB.CreateBookmark(r.Context(), params)
	if errors.Is(err, sql.ErrNoRows) {
		response.WithError(w, http.StatusNotFound, "Post not found")
		return
	}
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Couldn't create bookmark")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleUnbookmarkPost removes a saved post for the authenticated user
func (h *Handler) HandleUnbookmarkPost(w http.ResponseWriter, r *http.Request, user database.User) {
	postID, err := uuid.Parse(chi.URLParam(r, "postID"))

	if err != nil {
		response.WithError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	err = h.DB.DeleteBookmark(r.Context(), database.DeleteBookmarkParams{
		UserID: user.ID,
		PostID: postID,
	})

	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Couldn't delete bookmark")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// HandleGetBookmarks returns a cursor-paginated list of the user's bookmarked posts.
func (h *Handler) HandleGetBookmarks(w http.ResponseWriter, r *http.Request, user database.User) {
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
	params := database.GetBookmarkedPostsParams{
		UserID: user.ID,
		Limit:  limit + 1,
	}
	if cursor != nil {
		params.CursorPublishedAt = sql.NullTime{Time: cursor.Timestamp, Valid: true}
		params.CursorID = uuid.NullUUID{UUID: cursor.ID, Valid: true}
	}

	posts, err := h.DB.GetBookmarkedPosts(r.Context(), params)

	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Couldn't fetch bookmarked posts")
		return
	}

	nextCursor := ""

	if len(posts) > int(limit) {
		posts = posts[:limit]
		last := posts[len(posts)-1]
		nextCursor = encodeCursor(last.PublishedAt, last.ID)
	}
	response.WithJSON(w, http.StatusOK, struct {
		Data       []model.Post `json:"data"`
		NextCursor string       `json:"next_cursor"`
	}{
		Data:       model.DatabasePostsToPosts(posts), // doing to this to control the shape of the response.
		NextCursor: nextCursor,
	})

}
