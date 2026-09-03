// handler_readiness.go
package main

import (
	"encoding/json"

	"log"
	"net/http"
	"time"

	"github.com/ShortOwl/RSSAggregator/internal/database"
	"github.com/google/uuid"
)

func (a *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// 1. Define what we expect from the JSON body
	type parameter struct {
		Name string `json:"name"`
	}

	// 2. Parse the JSON body

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("Error parsing the JSON body : %v", err)
		return
	}

	// 3. call the sqlc generate function
	user, err := a.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),       // generate new uuid
		CreatedAt: time.Now().UTC(), // current time stamp
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name, // from json body
	})

	if err != nil {
		respondWithError(w, 500, "Couldn't create user")
		return
	}

	// 4. Return created user as JSON
	respondWithJSON(w, 201, databaseUserToUser(user))

}

// handler for getting the user

func (a *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 201, databaseUserToUser(user))

}

func (a *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {

	posts, err := a.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(10),
	})

	if err != nil {
		log.Println("Error fetching posts for user with user_id :", user.ID, "err:", err)
	}

	respondWithJSON(w, 200, databasePoststoPosts(posts))
}
