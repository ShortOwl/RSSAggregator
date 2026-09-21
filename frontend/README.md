# Margin — RSS reader frontend

A Next.js App Router application that uses the existing Go REST API. Design tokens and component styling follow the root `DESIGN.md`.

## Run

Use Node.js 22.13 or newer.

```sh
cd frontend
npm install
cp .env.example .env.local
npm run dev
```

Set `NEXT_PUBLIC_API_URL` to the Go API including `/v1` (default `http://localhost:8080/v1`). Start the existing Go server and PostgreSQL separately. Open http://localhost:3000. For production, set this public URL before building, use HTTPS for both services, and run `npm run build` followed by `npm start`.

## Components

- `app/`: six routes, shared layouts, and query providers.
- `components/auth/`: login and registration, session restoration, expiry, and protected navigation.
- `components/posts/`: safe text excerpts, feed/search/unread filters, paginated lists, and reading actions.
- `components/feeds/`: directory, subscriptions, and add-then-follow recovery.
- `components/settings/`: read-only account details and browser-local compact reading preference.
- `components/ui/`: local shadcn/ui primitives built with Radix and Tailwind, restyled to the supplied design.
- `hooks/`: API queries and mutations; no backend business logic.
- `lib/api.ts`: Fetch transport, Bearer headers, backend errors, and 204 handling.

## API integration

Login and registration use `/login` and `/register`. The issued JWT is stored in browser local storage and sent as a Bearer token directly to the Go API. Sessions end at the existing 24-hour expiry or on a protected 401 response; there is no refresh endpoint. Signing out clears credentials and the query cache, including across tabs. This browser-persisted token is accessible to JavaScript; an HttpOnly-cookie deployment would require a separate server transport design.

Posts use `/posts` with `search`, `feed_id`, `unread`, `limit`, and the server's opaque `cursor`. Bookmarks use `/bookmarks` with pagination only. An empty `next_cursor` ends the list. Infinite scrolling includes a manual load/retry button. Search is debounced. Failed requests honor rate-limit delays through TanStack Query retries.

Read and bookmark mutations use PUT/DELETE `/posts/{id}/read` and `/posts/{id}/bookmark`. Because responses omit status flags, membership is reconciled through all bookmark and unread pages. Until reconciliation finishes, explicit actions are available without assuming state. Unread reconciliation only establishes status for followed feeds. Large collections require additional requests; pages are requested sequentially. Status lookups can be retried without disabling the reading list.

Feed names are resolved from the cursor-paginated `/feeds` directory. `/feed_follows` provides subscription IDs; unfollow uses that ID rather than the feed ID. Creation and following are separate operations. If following fails, retry uses the already-created feed. The scraper may take about ten minutes to populate a new feed.

Profile uses `GET /users`; account edits and password reset are omitted because the API does not provide them. The returned API key is not displayed. Article links allow only HTTP(S); descriptions render as text, never injected HTML.

## Verification

```sh
npm run typecheck
npm run lint
npm test
npx playwright install chromium
npm run test:e2e
npm run build
```

Browser tests intercept the existing API contract, covering desktop and mobile without modifying real accounts or requiring a seeded database. They do not replace a live-backend deployment smoke test.
