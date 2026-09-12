package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// The ordinary http.ResponseWriter lets us write a status,
// but does not provide a method to read it afterward. Our wrapper remembers it.

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(status int) {
	if w.status != 0 {
		return
	}
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// handle responses that write a body without setting a status.
func (w *responseWriter) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(data)
}

// Logger records one structured log entry for each HTTP request
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()
		wrapped := &responseWriter{ResponseWriter: w}

		next.ServeHTTP(wrapped, r)
		// We pass wrapped instead of the original w.
		// Consequently, when downstream code calls WriteHeader or Write,
		// our methods capture the status before forwarding the response.

		status := wrapped.status
		if status == 0 {
			status = http.StatusOK
		}
		slog.InfoContext(r.Context(), "request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", status,
			"latency", time.Since(startedAt),
			"request_id", GetRequestID(r.Context()),
			"remote_addr", r.RemoteAddr,
		)
	})
}
