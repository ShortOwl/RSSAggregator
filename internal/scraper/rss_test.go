package scraper

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ShortOwl/RSSAggregator/internal/database"
	"github.com/google/uuid"
)

func TestScrapeBatchCleanup(t *testing.T) {
	for _, tt := range []struct {
		name         string
		feedCount    int
		duplicate    bool
		invalidXML   bool
		listError    bool
		cleanupError bool
	}{
		{name: "cleanup after article insertion", feedCount: 1},
		{name: "ten workers clean up once", feedCount: 10},
		{name: "empty batch still cleans up"},
		{name: "duplicate still cleans up", feedCount: 1, duplicate: true},
		{name: "failed feed still cleans up", feedCount: 1, invalidXML: true},
		{name: "feed list failure skips cleanup", listError: true},
		{name: "cleanup failure is handled", feedCount: 1, cleanupError: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var logs bytes.Buffer
			previousOutput := log.Writer()
			log.SetOutput(&logs)
			defer log.SetOutput(previousOutput)

			conn, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			// Workers may issue queries in any order. Single-worker cases keep
			// strict ordering to catch cleanup running before the article insert.
			if tt.feedCount > 1 {
				mock.MatchExpectationsInOrder(false)
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.invalidXML {
					_, _ = w.Write([]byte("<invalid"))
					return
				}
				_, _ = w.Write([]byte("<rss><channel><item><title>Article</title><link>https://example.com/post</link><pubDate>Wed, 16 Sep 2026 12:00:00 +0000</pubDate></item></channel></rss>"))
			}))
			defer server.Close()
			columns := []string{"id", "created_at", "updated_at", "name", "url", "user_id", "last_fetched_at"}
			feeds := sqlmock.NewRows(columns)
			list := mock.ExpectQuery("SELECT .* from feeds").WithArgs(int32(10))
			if tt.listError {
				list.WillReturnError(errors.New("feed list unavailable"))
			} else {
				list.WillReturnRows(feeds)
			}
			now := time.Now()
			for i := 0; i < tt.feedCount; i++ {
				feedID, userID := uuid.New(), uuid.New()
				feeds.AddRow(feedID, now, now, "Feed", server.URL, userID, nil)
				mock.ExpectQuery("UPDATE feeds").WithArgs(feedID).WillReturnRows(
					sqlmock.NewRows(columns).AddRow(feedID, now, now, "Feed", server.URL, userID, now),
				)
				if !tt.invalidXML {
					insert := mock.ExpectQuery("INSERT INTO posts")
					if tt.duplicate {
						insert.WillReturnError(errors.New("duplicate key value violates unique constraint"))
					} else {
						insert.WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "title", "description", "published_at", "url", "feed_id"}).AddRow(uuid.New(), now, now, "Article", nil, now, "https://example.com/post", feedID))
					}
				}
			}
			if !tt.listError {
				cleanup := mock.ExpectExec("DELETE FROM posts").WithArgs(int32(30))
				if tt.cleanupError {
					cleanup.WillReturnError(errors.New("cleanup unavailable"))
				} else {
					cleanup.WillReturnResult(sqlmock.NewResult(0, 2))
				}
			}
			scrapeBatch(database.New(conn), 10, 30)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
			wantSuccess, wantFailure := 0, 0
			if !tt.listError {
				if tt.cleanupError {
					wantFailure = 1
				} else {
					wantSuccess = 1
				}
			}
			// Also catch extra cleanup calls whose SQL errors are logged rather
			// than returned, and verify the successful deleted-row count.
			if got := strings.Count(logs.String(), "Retention cleanup deleted 2 posts (retention: 30 days)"); got != wantSuccess {
				t.Fatalf("cleanup success logs = %d, want %d: %s", got, wantSuccess, logs.String())
			}
			if got := strings.Count(logs.String(), "Failed to delete expired posts:"); got != wantFailure {
				t.Fatalf("cleanup failure logs = %d, want %d: %s", got, wantFailure, logs.String())
			}
		})
	}
}

func TestFetchFeed(t *testing.T) {
	for _, tt := range []struct {
		name, fixture string
		wantErr       bool
	}{
		{"multiple items", "testdata/feed.xml", false},
		{"malformed XML", "testdata/invalid.xml", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(tt.fixture)
			if err != nil {
				t.Fatal(err)
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/rss+xml")
				_, _ = w.Write(data)
			}))
			defer server.Close()
			feed, err := fetchFeed(server.URL)
			if (err != nil) != tt.wantErr {
				t.Fatalf("fetchFeed() error=%v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if feed.Channel.Title != "Go & Backend" || feed.Channel.Link != "https://example.com" || feed.Channel.Description != "Local test feed" {
				t.Errorf("unexpected channel: %+v", feed.Channel)
			}
			want := []RSSItem{
				{Title: "First post", Link: "https://example.com/first", Description: "<p>Learn Go & RSS</p>", PubDate: "Wed, 16 Sep 2026 12:00:00 +0000"},
				{Title: "Second post", Link: "https://example.com/second", Description: "", PubDate: "Tue, 15 Sep 2026 12:00:00 +0000"},
			}
			if len(feed.Channel.Item) != len(want) {
				t.Fatalf("got %d items, want %d", len(feed.Channel.Item), len(want))
			}
			for i, item := range feed.Channel.Item {
				if item != want[i] {
					t.Errorf("item %d=%+v, want %+v", i, item, want[i])
				}
			}
		})
	}
}

func TestFetchFeedInvalidURL(t *testing.T) {
	if _, err := fetchFeed("://invalid"); err == nil {
		t.Error("invalid URL should fail")
	}
}

func TestParsePublishedAt(t *testing.T) {
	for _, tt := range []struct {
		name  string
		value string
		want  time.Time
	}{
		{
			name:  "numeric offset",
			value: "Sat, 19 Sep 2026 23:01:57 +0000",
			want:  time.Date(2026, time.September, 19, 23, 1, 57, 0, time.UTC),
		},
		{
			name:  "GMT zone",
			value: "Sat, 19 Sep 2026 23:01:57 GMT",
			want:  time.Date(2026, time.September, 19, 23, 1, 57, 0, time.UTC),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePublishedAt(tt.value)
			if err != nil {
				t.Fatalf("parsePublishedAt(%q): %v", tt.value, err)
			}
			if !got.Equal(tt.want) {
				t.Fatalf("parsePublishedAt(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}

	if _, err := parsePublishedAt("not a date"); err == nil {
		t.Fatal("parsePublishedAt should reject an invalid date")
	}
}
