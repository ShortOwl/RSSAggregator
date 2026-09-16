package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	// Zero refill makes burst exhaustion deterministic without sleeping.
	limiter := NewRateLimiter(0, 2)
	calls := 0
	handler := limiter.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusNoContent)
	}))
	for _, tt := range []struct {
		name, address, authorization string
		want                         int
	}{
		{"first IP request", "192.0.2.1:1000", "", 204},
		{"same IP different port", "192.0.2.1:2000", "", 204},
		{"IP exhausted", "192.0.2.1:3000", "", 429},
		{"different IP", "192.0.2.2:1000", "", 204},
		{"first API key request", "192.0.2.1:1000", "ApiKey first", 204},
		{"same key different IP", "192.0.2.2:1000", "apikey first", 204},
		{"key exhausted", "192.0.2.3:1000", "ApiKey first", 429},
		{"different key", "192.0.2.1:1000", "ApiKey second", 204},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			r.RemoteAddr = tt.address
			r.Header.Set("Authorization", tt.authorization)
			w := httptest.NewRecorder()
			before := calls
			handler.ServeHTTP(w, r)
			if w.Code != tt.want {
				t.Fatalf("status=%d, want %d", w.Code, tt.want)
			}
			if tt.want == 429 {
				if calls != before {
					t.Error("rejected request reached handler")
				}
				if w.Header().Get("Retry-After") != "1" {
					t.Error("missing Retry-After")
				}
				var body struct {
					Error string `json:"error"`
				}
				if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil || body.Error != "rate limit exceeded" {
					t.Errorf("unexpected error body: %s", w.Body.String())
				}
			} else if calls != before+1 {
				t.Error("allowed request did not reach handler")
			}
		})
	}
}

func TestRateLimiterRemovesInactiveClients(t *testing.T) {
	limiter := NewRateLimiter(0, 1)
	if !limiter.allow("inactive") || !limiter.allow("active") {
		t.Fatal("new clients should be allowed")
	}
	// Age entries directly instead of waiting for the cleanup interval.
	past := time.Now().Add(-2 * limiterEntryTTL)
	limiter.clients["inactive"].lastSeen = past
	limiter.lastCleanup = past
	if !limiter.allow("new") {
		t.Fatal("new client rejected")
	}
	if _, exists := limiter.clients["inactive"]; exists {
		t.Error("inactive client was not removed")
	}
	if limiter.allow("active") {
		t.Error("active client's exhausted bucket was reset")
	}
}
