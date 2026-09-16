package middleware

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ShortOwl/RSSAggregator/internal/auth"
	"github.com/ShortOwl/RSSAggregator/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestWithAuthRejectsCredentials(t *testing.T) {
	expired, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: auth.TokenIssuer, Subject: uuid.NewString(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
	}).SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct{ name, header, message string }{
		{"missing", "", "Authentication credentials are required"},
		{"missing token", "Bearer", "Authentication credentials are required"},
		{"missing API key", "ApiKey", "Authentication credentials are required"},
		{"extra fields", "Bearer one two", "Authentication credentials are required"},
		{"unsupported", "Basic abc", "Unsupported authentication scheme"},
		{"malformed JWT", "Bearer bad", "token is expired or invalid"},
		{"expired JWT", "Bearer " + expired, "token is expired or invalid"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// These requests must fail before any database lookup.
			handler := WithAuth(nil, "test-secret", func(http.ResponseWriter, *http.Request, database.User) {
				t.Error("rejected credentials reached protected handler")
			})
			r := httptest.NewRequest("GET", "/users", nil)
			r.Header.Set("Authorization", tt.header)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)
			var body struct {
				Error string `json:"error"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if w.Code != 401 || body.Error != tt.message {
				t.Fatalf("status=%d, body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestWithAuthUserLookup(t *testing.T) {
	id := uuid.New()
	token, err := auth.GenerateJWT(id, "test-secret")
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name, header, query string
		lookupErr           error
		want                int
	}{
		{"Bearer success", "bEaReR " + token, "GetUserByID", nil, 204},
		{"API key success", "aPiKeY test-key", "GetUserByAPIKey", nil, 204},
		{"user missing", "Bearer " + token, "GetUserByID", sql.ErrNoRows, 401},
		{"invalid API key", "ApiKey test-key", "GetUserByAPIKey", sql.ErrNoRows, 401},
		{"database failure", "Bearer " + token, "GetUserByID", errors.New("database unavailable"), 500},
	} {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			query := mock.ExpectQuery(tt.query)
			if tt.query == "GetUserByID" {
				query.WithArgs(id)
			} else {
				query.WithArgs("test-key")
			}
			if tt.lookupErr != nil {
				query.WillReturnError(tt.lookupErr)
			} else {
				query.WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "api_key", "email", "password_hash"}).
					AddRow(id, time.Now(), time.Now(), "Reader", "test-key", "reader@example.com", "hash"))
			}
			called := false
			handler := WithAuth(database.New(db), "test-secret", func(w http.ResponseWriter, r *http.Request, user database.User) {
				called = true
				if user.ID != id || user.Name != "Reader" {
					t.Errorf("unexpected user: %+v", user)
				}
				w.WriteHeader(http.StatusNoContent)
			})
			r := httptest.NewRequest("GET", "/users", nil)
			r.Header.Set("Authorization", tt.header)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)
			if w.Code != tt.want || called != (tt.want == 204) {
				t.Errorf("status=%d, called=%v; want %d", w.Code, called, tt.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}
