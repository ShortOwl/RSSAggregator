package handler

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ShortOwl/RSSAggregator/internal/database"
	"github.com/ShortOwl/RSSAggregator/internal/model"
	"github.com/google/uuid"
)

func TestHandleGetFeeds(t *testing.T) {
	timestamp := time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC)
	ids := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}
	for _, tt := range []struct {
		name                string
		count               int
		withCursor, dbError bool
	}{
		{"empty page", 0, false, false},
		{"last page", 1, false, false},
		{"exact limit", 2, false, false},
		{"has next page", 3, false, false},
		{"incoming cursor", 1, true, false},
		{"database failure", 0, false, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			target := "/feeds?limit=2&search=%20Go%20"
			query := mock.ExpectQuery("GetFeeds")
			if tt.withCursor {
				target += "&cursor=" + url.QueryEscape(encodeCursor(timestamp, ids[0]))
				query.WithArgs("Go", timestamp, ids[0], 3)
			} else {
				query.WithArgs("Go", nil, nil, 3)
			}
			if tt.dbError {
				query.WillReturnError(errors.New("database unavailable"))
			} else {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "url", "user_id", "last_fetched_at"})
				for i := 0; i < tt.count; i++ {
					rows.AddRow(ids[i], timestamp.Add(-time.Duration(i)*time.Hour), timestamp, "Go feed", "https://example.com/rss", uuid.New(), nil)
				}
				query.WillReturnRows(rows).RowsWillBeClosed()
			}
			w := httptest.NewRecorder()
			New(database.New(db), "test-secret").HandleGetFeeds(w, httptest.NewRequest("GET", target, nil))
			if tt.dbError {
				if w.Code != 500 || w.Body.String() != `{"error":"Error fetching feeds"}` {
					t.Errorf("unexpected failure: %d %s", w.Code, w.Body.String())
				}
			} else {
				if w.Code != 200 {
					t.Fatalf("status=%d; %s", w.Code, w.Body.String())
				}
				var body struct {
					Data       []model.Feed `json:"data"`
					NextCursor string       `json:"next_cursor"`
				}
				if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
					t.Fatal(err)
				}
				wantCount := tt.count
				if wantCount > 2 {
					wantCount = 2
				}
				if body.Data == nil || len(body.Data) != wantCount {
					t.Fatalf("data=%v, want %d items (empty array for no results)", body.Data, wantCount)
				}
				for i, feed := range body.Data {
					if feed.ID != ids[i] || feed.Name != "Go feed" {
						t.Errorf("unexpected feed at index %d: %+v", i, feed)
					}
				}
				if tt.count > 2 {
					cursor, err := decodeCursor(body.NextCursor)
					if err != nil || cursor == nil {
						t.Fatalf("missing/invalid next cursor: %q", body.NextCursor)
					}
					if cursor.ID != ids[1] || !cursor.Timestamp.Equal(timestamp.Add(-time.Hour)) {
						t.Errorf("cursor must reference last returned feed: %+v", cursor)
					}
				} else if body.NextCursor != "" {
					t.Errorf("unexpected next cursor: %q", body.NextCursor)
				}
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}
