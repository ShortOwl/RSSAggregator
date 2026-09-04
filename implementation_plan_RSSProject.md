# RSS Aggregator → Production-Grade Go Backend

Transform this from a tutorial-level RSS scraper into a **portfolio-grade backend service** showcasing production Go engineering skills.

## Why This Plan Matters for Your CV

Interviewers look for **evidence of production thinking**, not just CRUD endpoints. This plan adds exactly the features that signal seniority: proper auth, observability, resilience, testing discipline, and deployment readiness. After this, you'll be able to describe the project as:

> *"A production-grade RSS aggregation API built in Go — featuring JWT authentication, cursor-based pagination, full-text search, request rate limiting, structured logging, comprehensive test coverage, Swagger documentation, Dockerized deployment, and a background scraping pipeline with configurable concurrency."*

---

## Proposed Changes

### Phase 1 — Project Structure & Code Quality

Restructure the flat layout into an idiomatic Go project structure. This alone signals engineering maturity.

#### [MODIFY] Project root — reorganize into packages

```
├── cmd/
│   └── server/
│       └── main.go            ← entry point
├── internal/
│   ├── auth/                  ← JWT + API key auth
│   ├── config/                ← centralized config (env parsing)
│   ├── database/              ← sqlc generated code
│   ├── handler/               ← HTTP handlers (users, feeds, posts, etc.)
│   ├── middleware/             ← auth, rate-limit, logging, recovery
│   ├── model/                 ← API response models + converters
│   ├── scraper/               ← RSS scraping engine
│   └── validator/             ← input validation
├── sql/
│   ├── queries/
│   └── schema/
├── docs/                      ← Swagger/OpenAPI spec
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── README.md
├── sqlc.yaml
├── go.mod
└── go.sum
```

#### [NEW] `Makefile`
Single-command build, test, lint, migrate, and run workflows.

#### [NEW] `README.md`
Professional README with architecture diagram, setup instructions, API overview, and tech stack. **This is the first thing recruiters see on GitHub.**

---

### Phase 2 — Authentication Upgrade (API Key → JWT)

Current auth is a raw API key in the `Authorization` header. Replace with **JWT (JSON Web Tokens)** — the industry standard.

#### [NEW] `sql/schema/007_user_password.sql`
Add `password_hash` column to users table.

#### [NEW] `sql/queries/auth.sql`
Queries for user registration and login lookups.

#### [MODIFY] `internal/auth/auth.go`
- Add `HashPassword()` and `CheckPasswordHash()` using `golang.org/x/crypto/bcrypt`
- Add `GenerateJWT()` and `ValidateJWT()` using `github.com/golang-jwt/jwt/v5`
- Keep backward-compatible API key auth as a fallback

#### [NEW] `internal/handler/handler_auth.go`
- `POST /v1/register` — register with name + email + password, returns JWT
- `POST /v1/login` — authenticate, returns JWT + refresh token

#### [MODIFY] `internal/middleware/middleware_auth.go`
Support both `Bearer <jwt>` and `apikey <key>` in `Authorization` header.

> [!IMPORTANT]
> **Decision needed**: Should we keep API key auth for backward compatibility, or fully replace it with JWT-only auth?

---

### Phase 3 — Pagination, Filtering & Search

Current endpoints return all results with no pagination — a red flag in any code review.

#### [NEW] `sql/schema/008_posts_search_index.sql`
Add GIN index on `posts.title` and `posts.description` for PostgreSQL full-text search.

#### [MODIFY] `sql/queries/posts.sql`
- Cursor-based pagination (using `published_at` + `id` as cursor)
- Full-text search: `WHERE to_tsvector('english', title || ' ' || COALESCE(description, '')) @@ plainto_tsquery('english', $1)`
- Filter by feed_id, date range

#### [MODIFY] `sql/queries/feeds.sql`
- Pagination for `GetFeeds`
- Search by name/URL

#### [MODIFY] Handlers
- Accept `?cursor=`, `?limit=`, `?search=`, `?feed_id=` query parameters
- Return pagination metadata in response: `{ data: [...], next_cursor: "..." }`

---

### Phase 4 — New Features (Bookmarks, Read Status, Categories)

These features make the project feel like a **real product**, not a tutorial exercise.

#### [NEW] `sql/schema/009_bookmarks.sql`
```sql
CREATE TABLE bookmarks (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    UNIQUE (user_id, post_id)
);
```

#### [NEW] `sql/schema/010_read_status.sql`
```sql
CREATE TABLE read_posts (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    read_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, post_id)
);
```

#### [NEW] `sql/schema/011_categories.sql`
```sql
CREATE TABLE categories (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (name, user_id)
);

CREATE TABLE feed_categories (
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (feed_id, category_id)
);
```

#### [NEW] Handlers
- `POST/DELETE /v1/bookmarks` — save/unsave posts
- `GET /v1/bookmarks` — get bookmarked posts (paginated)
- `POST /v1/posts/{postID}/read` — mark as read
- `GET /v1/posts?unread=true` — filter unread
- CRUD for `/v1/categories`
- `POST /v1/feeds/{feedID}/categories` — assign feed to category

---

### Phase 5 — Middleware Stack (Rate Limiting, Logging, Recovery)

This is where you demonstrate **production awareness**.

#### [NEW] `internal/middleware/rate_limiter.go`
Token-bucket rate limiter per API key / IP using `golang.org/x/time/rate`.

#### [NEW] `internal/middleware/logger.go`
Structured request logging with `log/slog` (Go stdlib):
- Method, path, status code, latency, request ID
- JSON output format for production, text for development

#### [NEW] `internal/middleware/recovery.go`
Panic recovery middleware that logs stack traces and returns 500.

#### [NEW] `internal/middleware/request_id.go`
Inject `X-Request-ID` header for tracing.

---

### Phase 6 — Configuration & Error Handling

#### [NEW] `internal/config/config.go`
Centralized config struct parsed from env vars with validation:
```go
type Config struct {
    Port            string
    DatabaseURL     string
    JWTSecret       string
    ScrapeInterval  time.Duration
    ScrapeConcurrency int
    RateLimitRPS    float64
    Environment     string // "development" | "production"
}
```

#### [MODIFY] Error responses
Standardize all errors to:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "human readable message",
    "details": { ... }
  }
}
```

---

### Phase 7 — Testing

Tests are **non-negotiable** for a CV project — they signal engineering discipline.

#### [NEW] `internal/auth/auth_test.go`
Unit tests for JWT generation/validation, password hashing.

#### [NEW] `internal/handler/*_test.go`
Integration tests using `httptest` for every endpoint.

#### [NEW] `internal/middleware/*_test.go`
Unit tests for rate limiter, auth middleware.

#### [NEW] `internal/scraper/scraper_test.go`
Test RSS parsing with fixture XML files.

---

### Phase 8 — API Documentation (Swagger/OpenAPI)

#### [NEW] `docs/swagger.yaml`
Full OpenAPI 3.0 spec documenting every endpoint, request/response schema, auth methods.

#### [MODIFY] `cmd/server/main.go`
Serve Swagger UI at `/docs` using `github.com/swaggo/http-swagger`.

---

### Phase 9 — Docker & Deployment

#### [NEW] `Dockerfile`
Multi-stage build: compile in `golang:1.27-alpine`, run in `scratch`/`alpine` — tiny final image.

#### [NEW] `docker-compose.yml`
Full local dev stack: Go server + PostgreSQL + optional pgAdmin.

#### [NEW] `.github/workflows/ci.yml`
GitHub Actions CI pipeline: lint (`golangci-lint`), test, build.

---

## Open Questions

> [!IMPORTANT]
> 1. **JWT vs API key**: Keep both auth mechanisms, or replace API key entirely?
> 2. **Email on users**: Should we add an `email` field to users for the registration flow?
> 3. **Deployment target**: Any preference — Railway, Fly.io, AWS, or just Docker-ready?
> 4. **OPML import/export**: Want the ability to import/export feed lists in OPML format (common RSS standard)?

---

## Verification Plan

### Automated Tests
```bash
make test          # run all unit + integration tests
make lint          # golangci-lint
make build         # ensure clean compile
```

### Manual Verification
- Run `docker-compose up` and test all endpoints via curl / Postman
- Verify Swagger UI loads at `/docs`
- Test JWT auth flow end-to-end
- Verify rate limiter triggers on burst requests
- Check structured logs output in JSON format

---

## Summary of New Dependencies

| Package | Purpose |
|---------|---------|
| `golang.org/x/crypto` | bcrypt password hashing |
| `github.com/golang-jwt/jwt/v5` | JWT tokens |
| `golang.org/x/time/rate` | Rate limiting |
| `github.com/swaggo/http-swagger` | Swagger UI |

All other improvements use **Go stdlib only** (`log/slog`, `net/http/httptest`, `testing`).
