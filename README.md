# RSS Aggregator — Backend Server

A production-grade RSS feed aggregation API built in Go. Users can subscribe to RSS feeds, and the server automatically scrapes and stores new posts in the background.

## Tech Stack

| Technology | Purpose |
|---|---|
| **Go** | Core language |
| **PostgreSQL** | Database |
| **Chi** | HTTP router |
| **sqlc** | Type-safe SQL code generation |
| **Goose** | Database migrations |

## Architecture

```
cmd/server/          → Entry point — wires dependencies and starts the server
internal/
  ├── auth/          → API key extraction from HTTP headers
  ├── config/        → Centralized environment configuration
  ├── database/      → Auto-generated database layer (sqlc)
  ├── handler/       → HTTP request handlers (users, feeds, posts)
  ├── middleware/     → Authentication middleware
  ├── model/         → API response models and DB-to-API converters
  ├── response/      → JSON response helpers
  └── scraper/       → Background RSS feed scraping engine
sql/
  ├── schema/        → Database migration files (goose)
  └── queries/       → SQL queries (sqlc)
```

## API Endpoints

### Public
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/v1/healthz` | Health check |
| `POST` | `/v1/users` | Create a new user |
| `GET` | `/v1/feeds` | List all feeds |

### Authenticated (requires `Authorization: apikey <key>`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/v1/users` | Get current user |
| `POST` | `/v1/feeds` | Add a new RSS feed |
| `POST` | `/v1/feed_follows` | Follow a feed |
| `GET` | `/v1/feed_follows` | List followed feeds |
| `DELETE` | `/v1/feed_follows/{id}` | Unfollow a feed |
| `GET` | `/v1/posts` | Get posts from followed feeds |

## Getting Started

### Prerequisites
- Go 1.21+
- PostgreSQL
- [goose](https://github.com/pressly/goose) (migrations)
- [sqlc](https://sqlc.dev/) (code generation)

### Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/ShortOwl/RSSAggregator.git
   cd RSSAggregator
   ```

2. **Create a `.env` file**
   ```
   PORT=8080
   DB_URL=postgres://username:password@localhost:5432/rssagg?sslmode=disable
   ```

3. **Run database migrations**
   ```bash
   make migrate-up
   ```

4. **Start the server**
   ```bash
   make run
   ```

### Available Commands
```bash
make build        # Compile binary to bin/
make run          # Run the server
make test         # Run all tests
make lint         # Run linter
make sqlc         # Regenerate database code
make migrate-up   # Apply database migrations
make migrate-down # Rollback last migration
```

## Background Scraper

The server runs a background goroutine that periodically fetches RSS feeds and stores new posts. It processes feeds concurrently using a configurable number of goroutines and deduplicates posts by URL.
