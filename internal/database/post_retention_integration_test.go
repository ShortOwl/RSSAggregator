package database

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// Run against a disposable PostgreSQL database, never the production database.
// Each run creates an isolated schema and removes it after verification.
func TestPostRetentionPostgres(t *testing.T) {
	dsn := os.Getenv("POST_RETENTION_TEST_DB_URL")
	if dsn == "" {
		t.Skip("set POST_RETENTION_TEST_DB_URL to run PostgreSQL integration test")
	}
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	conn.SetMaxOpenConns(1)
	ctx := context.Background()
	exec := func(query string, args ...any) {
		t.Helper()
		if _, err := conn.ExecContext(ctx, query, args...); err != nil {
			t.Fatal(err)
		}
	}
	schema := "retention_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	exec("CREATE SCHEMA " + schema)
	defer conn.ExecContext(ctx, "DROP SCHEMA "+schema+" CASCADE")
	exec("SET search_path TO " + schema)
	exec("SET TIME ZONE 'UTC'")
	files, err := filepath.Glob("../../sql/schema/*.sql")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		exec(strings.Split(string(data), "-- +goose Down")[0])
	}
	userID, feedID := uuid.New(), uuid.New()
	exec("INSERT INTO users (id, created_at, updated_at, name, email, password_hash) VALUES ($1, NOW(), NOW(), 'Reader', 'reader@example.com', 'test')", userID)
	exec("INSERT INTO feeds (id, created_at, updated_at, name, url, user_id) VALUES ($1, NOW(), NOW(), 'Feed', 'https://example.com/feed', $2)", feedID, userID)
	oldID := uuid.New()
	exec("INSERT INTO posts (id, created_at, updated_at, title, published_at, url, feed_id) VALUES ($1, NOW(), NOW(), 'Old', NOW() - INTERVAL '91 days', 'https://example.com/old', $2)", oldID, feedID)
	exec("INSERT INTO bookmarks (user_id, post_id) VALUES ($1, $2)", userID, oldID)
	exec("INSERT INTO read_posts (user_id, post_id) VALUES ($1, $2)", userID, oldID)
	newPost := func(url string, days int) CreatePostParams {
		now := time.Now().UTC()
		return CreatePostParams{ID: uuid.New(), CreatedAt: now, UpdatedAt: now, Title: url, PublishedAt: now.AddDate(0, 0, -days), Url: url, FeedID: feedID}
	}
	db := New(conn)
	recent := newPost("https://example.com/recent", 1)
	older := newPost("https://example.com/older", 60)
	for _, post := range []CreatePostParams{recent, older, newPost("https://example.com/ancient", 120)} {
		if _, err := db.CreatePost(ctx, post); err != nil {
			t.Fatal(err)
		}
	}
	deleted, err := db.DeleteExpiredPosts(ctx, 90)
	if err != nil || deleted != 2 {
		t.Fatalf("sync: deleted=%d err=%v", deleted, err)
	}
	count := func(table string, want int) {
		t.Helper()
		var got int
		if err := conn.QueryRowContext(ctx, "SELECT count(*) FROM "+table).Scan(&got); err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("%s count=%d want=%d", table, got, want)
		}
	}
	count("posts", 2)
	count("bookmarks", 0)
	count("read_posts", 0)
	// A duplicate does not prevent subsequent inserts or cleanup.
	if _, err := db.CreatePost(ctx, recent); err == nil {
		t.Fatal("expected duplicate insert to fail")
	}
	if _, err := db.CreatePost(ctx, newPost("https://example.com/another", 1)); err != nil {
		t.Fatal(err)
	}
	deleted, err = db.DeleteExpiredPosts(ctx, 30)
	if err != nil || deleted != 1 {
		t.Fatalf("duplicate/custom retention: deleted=%d err=%v", deleted, err)
	}
	count("posts", 2)
	deleted, err = db.DeleteExpiredPosts(ctx, 30)
	if err != nil || deleted != 0 {
		t.Fatalf("repeated cleanup: deleted=%d err=%v", deleted, err)
	}
	// A deletion failure must leave successfully inserted articles intact.
	exec("INSERT INTO posts (id, created_at, updated_at, title, published_at, url, feed_id) VALUES ($1, NOW(), NOW(), 'Old', NOW() - INTERVAL '91 days', 'https://example.com/old', $2)", oldID, feedID)
	exec("CREATE FUNCTION reject_delete() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'test cleanup failure'; END $$")
	exec("CREATE TRIGGER reject_delete BEFORE DELETE ON posts FOR EACH ROW EXECUTE FUNCTION reject_delete()")
	if _, err := db.CreatePost(ctx, newPost("https://example.com/preserved", 0)); err != nil {
		t.Fatal(err)
	}
	deleted, err = db.DeleteExpiredPosts(ctx, 90)
	if err == nil || deleted != 0 {
		t.Fatalf("expected cleanup failure, got deleted=%d err=%v", deleted, err)
	}
	count("posts", 4)
	exec("DROP TRIGGER reject_delete ON posts")
	// Use one transaction so NOW() is identical for the strict cutoff check.
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, "INSERT INTO posts (id, created_at, updated_at, title, published_at, url, feed_id) VALUES ($1, NOW(), NOW(), 'Boundary', NOW() - INTERVAL '90 days', 'https://example.com/boundary', $2)", uuid.New(), feedID)
	if err != nil {
		t.Fatal(err)
	}
	deleted, err = db.WithTx(tx).DeleteExpiredPosts(ctx, 90)
	if err != nil || deleted != 1 {
		t.Fatalf("strict boundary: deleted=%d err=%v", deleted, err)
	}
	var remains bool
	if err := tx.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM posts WHERE url = 'https://example.com/boundary')").Scan(&remains); err != nil || !remains {
		t.Fatalf("cutoff boundary removed: remains=%v err=%v", remains, err)
	}
}
