// handler_readiness.go
package main

import (
	"encoding/json"
	"fmt"

	"log"
	"net/http"
	"time"

	"github.com/ShortOwl/RSSAggregator/internal/database"
	"github.com/google/uuid"
)

// handler for getting the user

func (a *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {

	type parameter struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("Error parsing the JSON body : %v", err)
		return
	}

	feed, err := a.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      user.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})

	respondWithJSON(w, 201, databaseFeedtoFeed(feed))
}

// Handler to get all the feeds

func (a *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := a.DB.GetFeeds(r.Context())

	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("error fetching feeds : ", err))
	}

	respondWithJSON(w, 201, databaseFeedstoFeeds(feeds))
}
