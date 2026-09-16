package middleware

import (
	"context"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestID(t *testing.T) {
	for _, tt := range []struct{ name, header, want string }{
		{"missing", "", ""}, {"blank", "   ", ""},
		{"preserved", "client-id", "client-id"}, {"trimmed", "  client-id  ", "client-id"},
		{"maximum length", strings.Repeat("a", 128), strings.Repeat("a", 128)},
		{"too long", strings.Repeat("a", 129), ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var contextID string
			handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				contextID = GetRequestID(r.Context())
				w.WriteHeader(http.StatusNoContent)
			}))
			r := httptest.NewRequest("GET", "/", nil)
			r.Header.Set("X-Request-ID", tt.header)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)
			got := w.Header().Get("X-Request-ID")
			if w.Code != http.StatusNoContent || got == "" || contextID != got {
				t.Fatalf("status=%d, response ID=%q, context ID=%q", w.Code, got, contextID)
			}
			if tt.want != "" {
				if got != tt.want {
					t.Errorf("ID=%q, want %q", got, tt.want)
				}
			} else {
				decoded, err := hex.DecodeString(got)
				if err != nil || len(decoded) != 16 {
					t.Errorf("invalid generated ID %q", got)
				}
			}
			if GetRequestID(r.Context()) != "" {
				t.Error("original request context was modified")
			}
		})
	}
	if got := GetRequestID(context.Background()); got != "" {
		t.Errorf("missing ID=%q", got)
	}
}
