// Package handler contains all HTTP request handlers for the RSS Aggregator API.
package handler

import "github.com/ShortOwl/RSSAggregator/internal/database"

// Handler holds shared dependencies (like the database connection) that
// all HTTP handlers need. By putting DB in a struct, we avoid global
// variables and make testing easier — you can create a Handler with a
// mock database for tests.
type Handler struct {
	DB *database.Queries
}

// New creates a new Handler with the given database queries.
func New(db *database.Queries) *Handler {
	return &Handler{DB: db}
}
