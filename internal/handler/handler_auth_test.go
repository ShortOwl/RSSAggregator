package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ShortOwl/RSSAggregator/internal/auth"
	"github.com/ShortOwl/RSSAggregator/internal/database"
	"github.com/google/uuid"
)

func TestHandleLogin(t *testing.T) {
	id := uuid.New()
	hash, err := auth.HashPassword("correct-password")
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name, password string
		lookupErr      error
		want           int
		message        string
	}{
		{"success", "correct-password", nil, 200, ""},
		{"wrong password", "wrong-password", nil, 401, "Invalid email or password"},
		{"unknown email", "correct-password", sql.ErrNoRows, 401, "Invalid email or password"},
		{"database failure", "correct-password", errors.New("database unavailable"), 500, "Couldn't log in"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			query := mock.ExpectQuery("GetUserByEmail").WithArgs("reader@example.com")
			if tt.lookupErr != nil {
				query.WillReturnError(tt.lookupErr)
			} else {
				query.WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "api_key", "email", "password_hash"}).
					AddRow(id, time.Now(), time.Now(), "Reader", "key", "reader@example.com", hash))
			}
			// The query expectation also verifies email trimming and lowercasing.
			body := `{"email":"  READER@example.com  ","password":"` + tt.password + `"}`
			w := httptest.NewRecorder()
			New(database.New(db), "test-secret").HandleLogin(w, httptest.NewRequest("POST", "/login", strings.NewReader(body)))
			if w.Code != tt.want {
				t.Fatalf("status=%d, want %d; %s", w.Code, tt.want, w.Body.String())
			}
			if tt.want == 200 {
				var response authResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatal(err)
				}
				got, err := auth.ValidateJWT(response.Token, "test-secret")
				if err != nil || got != id || response.TokenType != "Bearer" {
					t.Errorf("invalid login token: ID=%v, type=%q, error=%v", got, response.TokenType, err)
				}
			} else {
				var response struct {
					Error string `json:"error"`
				}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatal(err)
				}
				if response.Error != tt.message {
					t.Errorf("error=%q, want %q", response.Error, tt.message)
				}
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}
