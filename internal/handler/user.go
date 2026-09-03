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

// HandleCreateUser creates a new user account.
// This is a public endpoint — no authentication required.
func (h *Handler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)
	if err != nil {
		response.WithError(w, 400, "Error parsing JSON body")
		return
	}

	user, err := h.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		response.WithError(w, 500, "Couldn't create user")
		return
	}

	response.WithJSON(w, 201, model.DatabaseUserToUser(user))
}

// HandleGetUser returns the currently authenticated user's info.
func (h *Handler) HandleGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	response.WithJSON(w, 200, model.DatabaseUserToUser(user))
}
