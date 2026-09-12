package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
)

const requestIDHeader = "X-Request-ID" // aa HTTP header che, used to send tracking no

type requestIDContextKey struct{}

// RequestID adds a request ID to the request context and response headers.
// The router supplies next, and RequestID returns a handler that adds request-ID behavior around it.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := strings.TrimSpace(r.Header.Get(requestIDHeader))

		if requestID == "" || len(requestID) > 128 {
			requestID = newRequestID()
		}

		w.Header().Set(requestIDHeader, requestID) // write request ID in response.
		// store the id for server's other middleware
		ctx := context.WithValue(r.Context(), requestIDContextKey{}, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
		// Without this call, the processing will stop here and the route handler will never execute.
	})
}

// GetRequestID returns the request ID stored in the request context.
func GetRequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDContextKey{}).(string)
	return requestID
}

func newRequestID() string {
	var id [16]byte
	if _, err := rand.Read(id[:]); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(id[:])
}
