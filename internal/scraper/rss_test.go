package scraper

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

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
