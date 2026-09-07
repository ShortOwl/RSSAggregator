package handler

// This file contains pagination logic shared by posts and feeds.
// here I create some helper functions which are used in pagination logic.
import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const (
	defaultPageLimit = 10
	maxPageLimit     = 100
)

type paginationCursor struct {
	Timestamp time.Time `json:"timestamp"`
	ID        uuid.UUID `json:"id"`
}

// request mathi limit ni value kadhva mate.
func parseLimit(r *http.Request) (int32, error) {

	value := r.URL.Query().Get("limit")

	if value == "" {
		// set to default value.
		return defaultPageLimit, nil
	}

	limit, err := strconv.Atoi(value)
	if err != nil || limit < 1 || limit > maxPageLimit {
		return 0, fmt.Errorf("Limit must be between 1 and %d", maxPageLimit)
	}

	return int32(limit), nil
}

// request mathi cursor ne decode krva mate.
func parseCursor(r *http.Request) (*paginationCursor, error) {
	value := r.URL.Query().Get("cursor")
	return decodeCursor(value)
}
func decodeCursor(value string) (*paginationCursor, error) {
	if value == "" {
		return nil, nil // no cursor is passed in URL.
	}

	data, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {

		return nil, fmt.Errorf("invalid cursor")
	}

	var cursor paginationCursor

	err = json.Unmarshal(data, &cursor)
	if err != nil || cursor.ID == uuid.Nil || cursor.Timestamp.IsZero() {
		return nil, fmt.Errorf("invalid cursor")
	}

	return &cursor, nil
}

func encodeCursor(timestamp time.Time, id uuid.UUID) string {
	data, _ := json.Marshal(paginationCursor{Timestamp: timestamp, ID: id})

	return base64.RawURLEncoding.EncodeToString(data)
}
