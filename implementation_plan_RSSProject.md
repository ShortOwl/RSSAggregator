# RSS Aggregator → Production-Grade Full-Stack Platform

Build a production-quality RSS aggregation platform that demonstrates strong **Go backend engineering** and a polished **Next.js frontend**. The backend is implemented and understood manually; the frontend is AI-assisted while every generated component is reviewed and understood.

---

# Project Philosophy

* **Backend:** Handwritten, fully understood, beginner-friendly Go code.
* **Frontend:** AI-assisted development with manual review of every component.
* **Priority:** Readability over clever abstractions.
* **Goal:** Build a product that is genuinely usable, not just a CRUD portfolio project.

---

# Final Roadmap

| Phase                                      | Status      |
| ------------------------------------------ | ----------- |
| Phase 1 — Project Structure & Code Quality | ✅ Completed |
| Phase 2 — JWT Authentication               | ✅ Completed |
| Phase 3 — Pagination, Search & Filtering   | ✅ Completed |
| Phase 4 — Bookmarks & Read Status          | ✅ Completed |
| Phase 5 — Middleware                       | ✅ Completed |
| Phase 6 — Configuration & Error Handling   | ⏭️ Skipped  |
| Phase 7 — Testing                          | ⏸️ Deferred |
| Phase 8 — Swagger Documentation            | ⏸️ Deferred |
| **Phase 9 — Next.js Frontend**             | 🔜 Next     |
| Phase 10 — Cloud Database & Deployment     | Pending     |
| Phase 11 — Docker & CI/CD                  | Pending     |

> **Why Phase 6 is skipped:** The project already contains a Config system and JSON response helpers. Graceful shutdown and operational hardening add complexity without improving learning at this stage.

---

# Completed Backend Features

## Authentication

* JWT authentication using bcrypt password hashing
* User registration & login
* Bearer JWT authentication
* Legacy API-key support for protected routes

## RSS Platform

* Create RSS feeds
* Follow / unfollow feeds
* Background RSS scraper
* Store articles in PostgreSQL

## Search & Pagination

* Cursor-based pagination
* PostgreSQL full-text search
* Feed filtering
* Search by feed name

## Reading Experience

* Bookmark articles
* Mark read / unread
* View unread-only posts
* Paginated bookmarks

## Middleware

* Authentication
* Request IDs
* Structured logging
* Panic recovery
* Rate limiting
* CORS

---

# Current Backend Architecture

```text
cmd/server
      │
      ▼
Chi Router
      │
      ▼
Middleware
      │
      ▼
Handlers
      │
      ▼
sqlc Queries
      │
      ▼
PostgreSQL
```

This architecture should **not** be refactored unless absolutely necessary.

---

# Phase 9 — AI-Assisted Next.js Frontend

## Objective

Build a beautiful, responsive frontend that consumes the existing Go API without changing backend endpoints.

## Tech Stack

* Next.js (App Router)
* React
* TypeScript
* Tailwind CSS
* TanStack Query
* shadcn/ui

## Folder Structure

```text
frontend/
├── app/
│   ├── login/
│   ├── register/
│   ├── feeds/
│   ├── bookmarks/
│   └── settings/
├── components/
│   ├── layout/
│   ├── feed/
│   ├── post/
│   └── ui/
├── lib/
│   ├── api.ts
│   ├── auth.ts
│   └── query.ts
└── types/
```

## Pages

### Authentication

* Login
* Register
* JWT session handling
* Auto redirect after login

### Home Feed

* Infinite scroll
* Cursor pagination
* Search
* Feed filter
* Bookmark
* Read / unread

### Feed Management

* Create feed
* Follow feed
* Unfollow feed
* Search feeds

### Bookmarks

* Paginated bookmarks
* Remove bookmark
* Reading progress

### Settings

* Profile
* Logout
* Session information

## AI Frontend Rules

* Generate **one page at a time**
* Components should remain small and readable
* Prefer simple React patterns
* Do not regenerate existing files unnecessarily
* Explain only component responsibilities and API integration

---

# Phase 10 — Cloud Database & Deployment

## Objective

Move application data from the local machine to the cloud.

## Infrastructure

| Component | Platform        |
| --------- | --------------- |
| Frontend  | Vercel          |
| Backend   | Railway         |
| Database  | Neon PostgreSQL |

## Result

```text
Next.js (Vercel)
        │
        ▼
Go API (Railway)
        │
        ▼
Cloud PostgreSQL (Neon)
```

The MacBook stores only source code while application data lives in the cloud.

---

# Phase 11 — Docker & CI/CD

Docker is intentionally placed **after deployment**, so it solves a real problem instead of becoming an isolated learning topic.

## Docker

* Multi-stage Dockerfile
* Lightweight runtime image
* Existing Go entry point

## Docker Compose

* Go backend
* PostgreSQL
* Persistent volume

## GitHub Actions

Automatically run:

```bash
go test ./...
go vet ./...
go build ./...
```

---

# Deferred Phases

## Phase 7 — Testing

Implement later:

* JWT unit tests
* Pagination tests
* Middleware tests
* RSS parser tests
* HTTP handler tests using `httptest`

## Phase 8 — Swagger

Implement later:

* `docs/swagger.yaml`
* Interactive Swagger UI
* Bearer authentication documentation
* Complete OpenAPI specification

---

# Final Resume Description

> **RSS Aggregator Platform** — A production-grade full-stack application built with Go, PostgreSQL, sqlc, JWT authentication, cursor-based pagination, full-text search, concurrent RSS scraping, and a modern Next.js frontend featuring bookmarks, reading state, and infinite scrolling. Deployed with cloud infrastructure and documented through OpenAPI.
