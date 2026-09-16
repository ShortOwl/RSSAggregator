package handler

import (
	"encoding/base64"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestParseLimit(t *testing.T) {
	for _, tt := range []struct {
		value   string
		want    int32
		wantErr bool
	}{
		{"", 10, false}, {"1", 1, false}, {"25", 25, false}, {"100", 100, false},
		{"0", 0, true}, {"-1", 0, true}, {"101", 0, true},
		{"abc", 0, true}, {"1.5", 0, true}, {"9999999999999999999999", 0, true},
	} {
		t.Run(tt.value, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/feeds?limit="+url.QueryEscape(tt.value), nil)
			got, err := parseLimit(r)

			if got != tt.want || (err != nil) != tt.wantErr {
				t.Fatalf("parseLimit() = %d, %v; want %d, wantErr %v", got, err, tt.want, tt.wantErr)
			}
		})
	}
}

func TestCursorRoundTrip(t *testing.T) {
	timestamp := time.Date(2026, 9, 16, 12, 30, 0, 123456000, time.UTC)
	id := uuid.New()
	encoded := encodeCursor(timestamp, id)
	cursor, err := decodeCursor(encoded)
	if err != nil || cursor == nil {
		t.Fatalf("decodeCursor() = %v, %v", cursor, err)
	}
	if cursor.ID != id || !cursor.Timestamp.Equal(timestamp) {
		t.Fatalf("cursor lost ID or timestamp: %+v", cursor)
	}
	r := httptest.NewRequest("GET", "/feeds?cursor="+url.QueryEscape(encoded), nil)
	parsed, err := parseCursor(r)
	if err != nil || parsed == nil || *parsed != *cursor {
		t.Fatalf("parseCursor() = %v, %v; want %v", parsed, err, cursor)
	}
}

func TestDecodeCursorInvalid(t *testing.T) {
	for _, tt := range []struct{ name, value string }{
		{"base64", "%%%"},
		{"JSON", base64.RawURLEncoding.EncodeToString([]byte("not-json"))},
		{"missing fields", base64.RawURLEncoding.EncodeToString([]byte("{}"))},
		{"nil ID", encodeCursor(time.Now(), uuid.Nil)},
		{"zero timestamp", encodeCursor(time.Time{}, uuid.New())},
		{"invalid UUID", base64.RawURLEncoding.EncodeToString([]byte(`{"timestamp":"2026-09-16T12:00:00Z","id":"bad"}`))},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if cursor, err := decodeCursor(tt.value); err == nil || cursor != nil {
				t.Fatalf("decodeCursor() = %v, %v; want nil and error", cursor, err)
			}
		})
	}
}

func TestParseCursorMissing(t *testing.T) {
	for _, target := range []string{"/feeds", "/feeds?cursor="} {
		cursor, err := parseCursor(httptest.NewRequest("GET", target, nil))
		if cursor != nil || err != nil {
			t.Fatalf("parseCursor(%q) = %v, %v; want nil, nil", target, cursor, err)
		}
	}
}
