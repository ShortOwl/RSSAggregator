// Package scraper implements a background RSS feed scraping engine.
// It periodically fetches RSS feeds from the database, parses them,
// and stores new posts.
package scraper

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/ShortOwl/RSSAggregator/internal/database"
	"github.com/google/uuid"
)

// Start begins the background RSS scraping loop.
// It fetches up to 'concurrency' feeds at each interval and processes
// them concurrently using goroutines.
func Start(db *database.Queries, concurrency int, interval time.Duration, retentionDays int32) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// The for loop with <-ticker.C blocks until the next tick.
	// The first iteration runs immediately (before the first tick).
	for ; ; <-ticker.C {
		scrapeBatch(db, concurrency, retentionDays)
	}
}

// scrapeBatch processes one batch, then cleans up once after all workers finish.
func scrapeBatch(db *database.Queries, concurrency int, retentionDays int32) {
	feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
	if err != nil {
		log.Println("error fetching feeds:", err)
		return
	}

	wg := &sync.WaitGroup{}
	for _, feed := range feeds {
		wg.Add(1)

		go func(feed database.Feed) {
			defer wg.Done()
			ScrapeFeed(db, feed)
		}(feed)
	}
	wg.Wait()

	// Delete the expired posts.
	deletedPosts, err := db.DeleteExpiredPosts(context.Background(), retentionDays)
	if err != nil {
		log.Println("Failed to delete expired posts:", err)
		return
	}
	log.Printf("Retention cleanup deleted %d posts (retention: %d days)", deletedPosts, retentionDays)
}

func ScrapeFeed(db *database.Queries, feed database.Feed) {

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched:", err)
		return
	}

	rssFeed, err := fetchFeed(feed.Url)
	if err != nil {
		log.Println("error fetching feed from url:", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		t, err := parsePublishedAt(item.PubDate)
		if err != nil {
			log.Printf("Couldn't parse date %v with err %v", item.PubDate, err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: t,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates") {
				continue
			}
			log.Println("failed to create post:", err)
		}
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}

// parsePublishedAt accepts the two RFC 1123 variants commonly emitted by RSS
// feeds: numeric offsets ("+0000") and named zones ("GMT").
func parsePublishedAt(value string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC1123Z, value)
	if err == nil {
		return parsed, nil
	}

	return time.Parse(time.RFC1123, value)
}
