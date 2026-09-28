package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/ShortOwl/RSSAggregator/internal/response"
	"golang.org/x/time/rate"
)

const limiterEntryTTL = 3 * time.Minute

type limiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
	// each client needs 2 things, 1) token bucket, 2) when the client made the last request
}

// RateLimiter applies an independent token bucket to each IP address.
type RateLimiter struct {
	mu          sync.Mutex
	clients     map[string]*limiterEntry
	limit       rate.Limit
	burst       int
	lastCleanup time.Time
}

// create a new rate limiter per client
func NewRateLimiter(limit rate.Limit, burst int) *RateLimiter {
	return &RateLimiter{
		clients:     make(map[string]*limiterEntry),
		limit:       limit,
		burst:       burst,
		lastCleanup: time.Now(),
	}
}

// Limit rejects requests when their client's token bucket is empty.
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := clientKey(r)
		allowed := rl.allow(key)
		if !allowed {
			w.Header().Set("Retry-After", "1")
			response.WithError(w, http.StatusTooManyRequests, "rate limit exceeded")
			return
		}

		next.ServeHTTP(w, r)

	})
}

func clientKey(r *http.Request) string {
	// Ignore the connection's port so that same IP address keeps the same bucket.
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "ip:" + r.RemoteAddr
	}
	return "ip:" + host
}

func (rl *RateLimiter) allow(key string) bool {
	now := time.Now()

	rl.mu.Lock()
	defer rl.mu.Unlock() // Requests run concurrently, so protect access to the shared clients map.

	// Check for inactive clients every 3 minutes, when a request arrives.
	if now.Sub(rl.lastCleanup) >= limiterEntryTTL {
		for key, entry := range rl.clients {
			if now.Sub(entry.lastSeen) >= limiterEntryTTL {
				delete(rl.clients, key)
			}
		}
		rl.lastCleanup = now
	}

	// a missing map entry returns a nil.Give new clients there own bucket.
	bucket := rl.clients[key]
	if bucket == nil {
		bucket = &limiterEntry{}
		bucket.limiter = rate.NewLimiter(rl.limit, rl.burst)
		rl.clients[key] = bucket
	}
	bucket.lastSeen = now

	// Allow consumes one token if available; otherwise it returns false.
	return bucket.limiter.Allow()
}
