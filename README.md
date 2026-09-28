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
  ├── auth/          → Password hashing and JWT operations
  ├── config/        → Centralized environment configuration
  ├── database/      → Auto-generated database layer (sqlc)
  ├── handler/       → HTTP request handlers (users, feeds, posts)
  ├── middleware/    → JWT authentication and rate limiting middleware
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
| `POST` | `/v1/register` | Register with name, email, and password |
| `POST` | `/v1/login` | Log in with email and password |
| `POST` | `/v1/forgot-password` | Email a six-digit recovery code, if the account exists |
| `POST` | `/v1/reset-password` | Set a new password with the email and recovery code |
| `GET` | `/v1/feeds` | List all feeds |

### Authenticated

Protected endpoints require `Authorization: Bearer <jwt>`. Legacy API-key
authentication is deprecated and no longer accepted. Existing database keys
remain stored for migration compatibility but are not returned by the API.

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/v1/users` | Get current user |
| `POST` | `/v1/feeds` | Add a new RSS feed |
| `POST` | `/v1/feed_follows` | Follow a feed |
| `GET` | `/v1/feed_follows` | List followed feeds |
| `DELETE` | `/v1/feed_follows/{id}` | Unfollow a feed |
| `GET` | `/v1/posts` | Get posts from followed feeds |
| `PUT` | `/v1/posts/{postID}/bookmark` | Bookmark a post |
| `DELETE` | `/v1/posts/{postID}/bookmark` | Remove a bookmark |
| `GET` | `/v1/bookmarks` | List bookmarked posts with cursor pagination |
| `PUT` | `/v1/posts/{postID}/read` | Mark a post as read |
| `DELETE` | `/v1/posts/{postID}/read` | Mark a post as unread |

Use `GET /v1/posts?unread=true` to return only unread posts.

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
   JWT_SECRET=replace-with-a-random-secret-at-least-32-characters-long
   SMTP_EMAIL=your-gmail-address@gmail.com
   SMTP_PASSWORD=your-google-app-password
   ```

   Password recovery sends codes through Gmail SMTP. Use a [Google App Password](https://support.google.com/accounts/answer/2461835) for `SMTP_PASSWORD`; the server requires both SMTP variables at startup. Keep `.env` out of version control.

3. **Run database migrations**

   Migration 007 makes email and password hashes mandatory. Remove any
   passwordless tutorial users before applying it.
   Migration 012 creates the password-recovery OTP table.

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

After all workers in each scraping batch finish, posts older than 50 days
(by `published_at`) are deleted globally using a separate cleanup query, unless
at least one user has bookmarked them. Bookmarked posts are retained regardless
of age. Removing the last bookmark makes an expired post eligible for the next
cleanup.
Set `POST_RETENTION_DAYS` to a positive integer to override the default of `50`.
Apply migration `011` before deploying to add the retention index. Cleanup logs
include the number of deleted posts. If cleanup fails, saved articles remain and
cleanup is retried after the next batch. Cleanup runs once per batch, even if
the batch is empty or individual feeds fail. If loading the feed list fails,
that batch skips cleanup. The existing article insertion
and duplicate handling remain unchanged.
Read markers for deleted posts are removed by existing cascading foreign keys.
Cleanup runs only while the scraper is running and syncing feeds.
