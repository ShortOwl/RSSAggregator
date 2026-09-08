# RSS Aggregator Database ER Diagram

This is the current design through Phase 4. Categories are intentionally not included.

```mermaid
erDiagram
    USERS ||--o{ FEEDS : creates
    USERS ||--o{ FEED_FOLLOWS : follows_through
    FEEDS ||--o{ FEED_FOLLOWS : is_followed_through
    FEEDS ||--o{ POSTS : publishes
    USERS ||--o{ BOOKMARKS : saves
    POSTS ||--o{ BOOKMARKS : is_saved_in
    USERS ||--o{ READ_POSTS : reads
    POSTS ||--o{ READ_POSTS : has_read_state_in

    USERS {
        uuid id PK
        timestamp created_at
        timestamp updated_at
        text name
        varchar api_key UK
        text email UK
        text password_hash
    }

    FEEDS {
        uuid id PK
        timestamp created_at
        timestamp updated_at
        text name
        text url UK
        uuid user_id FK
        timestamp last_fetched_at
    }

    FEED_FOLLOWS {
        uuid id PK
        timestamp created_at
        timestamp updated_at
        uuid user_id FK
        uuid feed_id FK
    }

    POSTS {
        uuid id PK
        timestamp created_at
        timestamp updated_at
        text title
        text description
        timestamp published_at
        text url UK
        uuid feed_id FK
    }

    BOOKMARKS {
        uuid user_id PK, FK
        uuid post_id PK, FK
        timestamp created_at
    }

    READ_POSTS {
        uuid user_id PK, FK
        uuid post_id PK, FK
        timestamp read_at
    }
```

## How to remember the design

- A user creates feeds and follows feeds.
- A feed publishes many posts.
- `feed_follows` connects users to the feeds they follow.
- `bookmarks` connects users to the posts they save.
- `read_posts` connects users to the posts they have read.
