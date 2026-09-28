package handler

import (
	"database/sql"
	"database/sql/driver"
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

func TestForgotPasswordDoesNotRevealUnknownEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	mock.ExpectQuery("GetUserByEmail").WithArgs("reader@example.com").WillReturnError(sql.ErrNoRows)
	w := httptest.NewRecorder()
	New(database.New(db), "test-secret", "", "").HandleForgotPassword(w,
		httptest.NewRequest("POST", "/forgot-password", strings.NewReader(`{"email":" READER@example.com "}`)))
	if w.Code != 200 || w.Body.String() != `{"message":"If that email is registered, a recovery code will be sent."}` {
		t.Fatalf("unexpected response: %d %s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

type passwordHashMatcher struct{}

func (passwordHashMatcher) Match(value driver.Value) bool {
	hash, ok := value.(string)
	return ok && auth.CheckPasswordHash("NewPassword123", hash) == nil
}

func TestResetPassword(t *testing.T) {
	for _, tt := range []struct {
		name       string
		storedHash string
		expiresAt  time.Time
		rows       int64
		wantStatus int
	}{
		{"valid code", auth.HashOTP("482913"), time.Now().Add(time.Minute), 1, 200},
		{"wrong code", auth.HashOTP("111111"), time.Now().Add(time.Minute), 0, 400},
		{"expired code", auth.HashOTP("482913"), time.Now().Add(-time.Minute), 0, 400},
		{"already used code", auth.HashOTP("482913"), time.Now().Add(time.Minute), 0, 400},
	} {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			id := uuid.New()
			mock.ExpectQuery("GetUserByEmail").WithArgs("reader@example.com").WillReturnRows(
				sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "api_key", "email", "password_hash"}).
					AddRow(id, time.Now(), time.Now(), "Reader", "key", "reader@example.com", "old-hash"))
			mock.ExpectQuery("GetPasswordResetByUser").WithArgs(id).WillReturnRows(
				sqlmock.NewRows([]string{"user_id", "otp_hash", "expires_at", "created_at"}).
					AddRow(id, tt.storedHash, tt.expiresAt, time.Now().Add(-time.Minute)))
			if tt.name == "valid code" || tt.name == "already used code" {
				mock.ExpectExec("CompletePasswordReset").WithArgs(
					id, auth.HashOTP("482913"), sqlmock.AnyArg(), passwordHashMatcher{}, sqlmock.AnyArg(),
				).WillReturnResult(sqlmock.NewResult(0, tt.rows))
			}
			w := httptest.NewRecorder()
			New(database.New(db), "test-secret", "", "").HandleResetPassword(w,
				httptest.NewRequest("POST", "/reset-password", strings.NewReader(
					`{"email":" READER@example.com ","otp":"482913","password":"NewPassword123"}`)))
			if w.Code != tt.wantStatus {
				t.Fatalf("status=%d, want %d; %s", w.Code, tt.wantStatus, w.Body.String())
			}
			if tt.wantStatus == 400 && w.Body.String() != `{"error":"Invalid or expired code"}` {
				t.Fatalf("reset exposed a different error: %s", w.Body.String())
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestResetPasswordUnknownEmailUsesGenericError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	mock.ExpectQuery("GetUserByEmail").WithArgs("reader@example.com").WillReturnError(sql.ErrNoRows)
	w := httptest.NewRecorder()
	New(database.New(db), "test-secret", "", "").HandleResetPassword(w,
		httptest.NewRequest("POST", "/reset-password", strings.NewReader(
			`{"email":"reader@example.com","otp":"482913","password":"NewPassword123"}`)))
	if w.Code != 400 || w.Body.String() != `{"error":"Invalid or expired code"}` {
		t.Fatalf("unexpected response: %d %s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestResetPasswordRejectsInvalidInputBeforeDatabase(t *testing.T) {
	for _, body := range []string{
		`{"email":"reader@example.com","otp":"12x456","password":"NewPassword123"}`,
		`{"email":"reader@example.com","otp":"482913","password":"short"}`,
	} {
		w := httptest.NewRecorder()
		New(nil, "test-secret", "", "").HandleResetPassword(w,
			httptest.NewRequest("POST", "/reset-password", strings.NewReader(body)))
		if w.Code != 400 {
			t.Fatalf("status=%d, want 400; %s", w.Code, w.Body.String())
		}
	}
}

func TestForgotPasswordDatabaseFailureKeepsGenericResponse(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	mock.ExpectQuery("GetUserByEmail").WillReturnError(errors.New("database unavailable"))
	w := httptest.NewRecorder()
	New(database.New(db), "test-secret", "", "").HandleForgotPassword(w,
		httptest.NewRequest("POST", "/forgot-password", strings.NewReader(`{"email":"reader@example.com"}`)))
	if w.Code != 200 || w.Body.String() != `{"message":"If that email is registered, a recovery code will be sent."}` {
		t.Fatalf("unexpected response: %d %s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
