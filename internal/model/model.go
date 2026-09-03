package model

import (
	"time"

	"github.com/ShortOwl/RSSAggregator/internal/database"
	"github.com/google/uuid"
)

// User is the API response model for a user.
type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
}

// DatabaseUserToUser converts a database user to an API response user.
func DatabaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		APIKey:    dbUser.ApiKey,
	}
}

// Feed is the API response model for a feed.
type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

// DatabaseFeedToFeed converts a database feed to an API response feed.
func DatabaseFeedToFeed(d database.Feed) Feed {
	return Feed{
		ID:        d.ID,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
		Name:      d.Name,
		Url:       d.Url,
		UserID:    d.UserID,
	}
}

// DatabaseFeedsToFeeds converts a slice of database feeds to API response feeds.
func DatabaseFeedsToFeeds(dbFeeds []database.Feed) []Feed {
	res := make([]Feed, 0, len(dbFeeds))
	for _, dbFeed := range dbFeeds {
		res = append(res, DatabaseFeedToFeed(dbFeed))
	}
	return res
}

// FeedFollow is the API response model for a feed follow.
type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

// DatabaseFeedFollowToFeedFollow converts a database feed follow to an API response.
func DatabaseFeedFollowToFeedFollow(dbFeedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        dbFeedFollow.ID,
		CreatedAt: dbFeedFollow.CreatedAt,
		UpdatedAt: dbFeedFollow.UpdatedAt,
		UserID:    dbFeedFollow.UserID,
		FeedID:    dbFeedFollow.FeedID,
	}
}

// DatabaseFeedFollowsToFeedFollows converts a slice of database feed follows.
func DatabaseFeedFollowsToFeedFollows(dbFeedFollows []database.FeedFollow) []FeedFollow {
	res := make([]FeedFollow, 0, len(dbFeedFollows))
	for _, ff := range dbFeedFollows {
		res = append(res, DatabaseFeedFollowToFeedFollow(ff))
	}
	return res
}

// Post is the API response model for a post.
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

// DatabasePostsToPosts converts a slice of database posts to API response posts.
func DatabasePostsToPosts(dbSlice []database.Post) []Post {
	res := make([]Post, 0, len(dbSlice))
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
