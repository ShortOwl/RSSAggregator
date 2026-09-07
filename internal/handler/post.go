package handler

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/ShortOwl/RSSAggregator/internal/database"
	"github.com/ShortOwl/RSSAggregator/internal/model"
	"github.com/ShortOwl/RSSAggregator/internal/response"
)

func (h *Handler) HandleGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	// parse limit first.
	// authenticated handler

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

	var feed_id uuid.NullUUID

	if value := r.URL.Query().Get("feed_id"); value != "" {

		id, err := uuid.Parse(value)
		if err != nil {
			response.WithError(w, http.StatusBadRequest, "invalid feed_id")
			return
		}
		feed_id = uuid.NullUUID{
			UUID:  id,
			Valid: true,
		}
	}

	// Building sqlc parameters to pass in query.
	params := database.GetPostsForUserParams{
		UserID: user.ID,
		FeedID: feed_id,
		Search: strings.TrimSpace(r.URL.Query().Get("search")),
		Limit:  limit + 1,
	}

	if cursor != nil {
		params.CursorPublishedAt = sql.NullTime{Time: cursor.Timestamp, Valid: true}
		params.CursorID = uuid.NullUUID{UUID: cursor.ID, Valid: true}
	}

	posts, err := h.DB.GetPostsForUser(r.Context(), params)

	if err != nil {
		log.Println("Error fetching posts for user:", user.ID, "err:", err)
		response.WithError(w, 500, "Couldn't fetch posts")
		return
	}

	nextCursor := ""
	// cursor ne apde Base64 string na form ma pass kriye che.
	if len(posts) > int(limit) {
		// we need to make cursor
		posts = posts[:limit]
		last := posts[len(posts)-1]
		nextCursor = encodeCursor(last.PublishedAt, last.ID)
	}
	type Payload struct {
		Data       []model.Post `json:"data"`
		NextCursor string       `json:"next_cursor"`
	}

	response.WithJSON(w, http.StatusOK, Payload{
		Data:       model.DatabasePostsToPosts(posts),
		NextCursor: nextCursor,
	})

}
