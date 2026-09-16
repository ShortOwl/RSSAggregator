package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ShortOwl/RSSAggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func TestPublicHealthHandlers(t *testing.T) {
	h := New(nil, "test-secret")
	for _, tt := range []struct {
		name    string
		handler http.HandlerFunc
		status  int
		body    string
	}{
		{"readiness", h.HandleReadiness, 200, `{}`},
		{"error", h.HandleError, 400, `{"error":"Something went wrong :("}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tt.handler(w, httptest.NewRequest("GET", "/", nil))
			if w.Code != tt.status || w.Body.String() != tt.body {
				t.Fatalf("status=%d, body=%s", w.Code, w.Body.String())
			}
			if w.Header().Get("Content-Type") != "application/json" {
				t.Error("missing JSON content type")
			}
		})
	}
}

func TestHandleGetUser(t *testing.T) {
	user := database.User{ID: uuid.New(), Name: "Reader", Email: "reader@example.com", ApiKey: "test-key", PasswordHash: "secret-password-hash"}
	w := httptest.NewRecorder()
	New(nil, "test-secret").HandleGetUser(w, httptest.NewRequest("GET", "/users", nil), user)
	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if w.Code != 200 || body["id"] != user.ID.String() || body["email"] != user.Email || body["name"] != user.Name {
		t.Fatalf("unexpected user response: %d %s", w.Code, w.Body.String())
	}
	if _, exists := body["password_hash"]; exists {
		t.Error("response exposes password_hash")
	}
	if strings.Contains(w.Body.String(), user.PasswordHash) {
		t.Error("response exposes password hash value")
	}
}

func TestHandlersRejectInvalidInput(t *testing.T) {
	// No database is needed: each input must fail before its first query.
	h := New(nil, "test-secret")
	user := database.User{ID: uuid.New()}
	router := chi.NewRouter()
	router.Post("/register", h.HandleRegister)
	router.Post("/login", h.HandleLogin)
	router.Get("/feeds", h.HandleGetFeeds)
	router.Get("/posts", func(w http.ResponseWriter, r *http.Request) { h.HandleGetPostsForUser(w, r, user) })
	router.Get("/bookmarks", func(w http.ResponseWriter, r *http.Request) { h.HandleGetBookmarks(w, r, user) })
	router.Put("/posts/{postID}/bookmark", func(w http.ResponseWriter, r *http.Request) { h.HandleBookmarkPost(w, r, user) })
	router.Put("/posts/{postID}/read", func(w http.ResponseWriter, r *http.Request) { h.HandleMarkPostAsRead(w, r, user) })
	for _, tt := range []struct{ name, method, target, body, message string }{
		{"register malformed JSON", "POST", "/register", `{`, "Couldn't decode request body"},
		{"login malformed JSON", "POST", "/login", `{`, "Couldn't decode request body"},
		{"missing name", "POST", "/register", `{"name":"  ","email":"reader@example.com","password":"password123"}`, "Name is required"},
		{"invalid email", "POST", "/register", `{"name":"Reader","email":"bad","password":"password123"}`, "A valid email is required"},
		{"short password", "POST", "/register", `{"name":"Reader","email":"reader@example.com","password":"short"}`, "Password must be at least 8 bytes"},
		{"long password", "POST", "/register", `{"name":"Reader","email":"reader@example.com","password":"` + strings.Repeat("a", 73) + `"}`, "Password must be at most 72 bytes"},
		{"feed limit", "GET", "/feeds?limit=0", "", "Limit must be between 1 and 100"},
		{"feed cursor", "GET", "/feeds?cursor=bad", "", "invalid cursor"},
		{"post feed filter", "GET", "/posts?feed_id=bad", "", "invalid feed_id"},
		{"post unread filter", "GET", "/posts?unread=bad", "", "invalid unread value"},
		{"bookmark cursor", "GET", "/bookmarks?cursor=bad", "", "invalid cursor"},
		{"bookmark post ID", "PUT", "/posts/bad/bookmark", "", "Invalid post ID"},
		{"read post ID", "PUT", "/posts/bad/read", "", "Invalid post ID"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, httptest.NewRequest(tt.method, tt.target, strings.NewReader(tt.body)))
			var body struct {
				Error string `json:"error"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if w.Code != 400 || body.Error != tt.message {
				t.Fatalf("status=%d, body=%s", w.Code, w.Body.String())
			}
		})
	}
}
