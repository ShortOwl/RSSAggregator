package main

import (
	"time"

	"github.com/ShortOwl/RSSAggregator/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		APIKey:    dbUser.ApiKey,
	}
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func databaseFeedtoFeed(d database.Feed) Feed {
	return Feed{
		ID:        d.ID,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
		Name:      d.Name,
		Url:       d.Url,
		UserID:    d.UserID,
	}
}

func databaseFeedstoFeeds(dbFeeds []database.Feed) []Feed {
	res := []Feed{}

	for _, dbFeed := range dbFeeds {
		res = append(res, Feed{
			ID:        dbFeed.ID,
			CreatedAt: dbFeed.CreatedAt,
			UpdatedAt: dbFeed.UpdatedAt,
			Name:      dbFeed.Name,
			Url:       dbFeed.Url,
			UserID:    dbFeed.UserID,
		})
	}
	return res
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func databaseFeedFollowToFeedFollow(dbFeedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        dbFeedFollow.ID,
		CreatedAt: dbFeedFollow.CreatedAt,
		UpdatedAt: dbFeedFollow.UpdatedAt,
		UserID:    dbFeedFollow.UserID,
		FeedID:    dbFeedFollow.FeedID,
	}
}

func databaseFeedFollowsToFeedFollows(dbFeedFollow []database.FeedFollow) []FeedFollow {
	res := []FeedFollow{}

	for _, dbFeedFollow := range dbFeedFollow {
		res = append(res, FeedFollow{
			ID:        dbFeedFollow.ID,
			CreatedAt: dbFeedFollow.CreatedAt,
			UpdatedAt: dbFeedFollow.UpdatedAt,
			UserID:    dbFeedFollow.UserID,
			FeedID:    dbFeedFollow.FeedID,
		})
	}
	return res
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Url         string    `json:"url"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func databasePoststoPosts(dbSlice []database.Post) []Post {

	res := []Post{}

	for _, it := range dbSlice {
		var description *string
		if it.Description.Valid {
			description = &it.Description.String
		}
		res = append(res, Post{
			ID:          it.ID,
			CreatedAt:   it.CreatedAt,
			UpdatedAt:   it.UpdatedAt,
			Title:       it.Title,
			Description: description,
			PublishedAt: it.PublishedAt,
			Url:         it.Url,
			FeedID:      it.FeedID,
		})
	}
	return res
}
