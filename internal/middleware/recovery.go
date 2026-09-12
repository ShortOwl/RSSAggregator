package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/ShortOwl/RSSAggregator/internal/response"
)

// Recovery converts panics into internal server error responses.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer recoverPanic(w, r) // recover panic when the next handler finishes or panics.
		next.ServeHTTP(w, r)
	})
}

func recoverPanic(w http.ResponseWriter, r *http.Request) {
	// recover must be called inside the defered function itself.
	panicValue := recover()
	if panicValue == nil {
		return
	}

	slog.ErrorContext(r.Context(), "panic recovered",
		"panic", panicValue,
		"request_id", GetRequestID(r.Context()),
		"stack", string(debug.Stack()),
	)
	response.WithError(w, 500, "Internal Server Error")
}
