package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/ShortOwl/RSSAggregator/internal/config"
	"github.com/ShortOwl/RSSAggregator/internal/database"
	"github.com/ShortOwl/RSSAggregator/internal/handler"
	"github.com/ShortOwl/RSSAggregator/internal/middleware"
	"github.com/ShortOwl/RSSAggregator/internal/scraper"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	_ "github.com/lib/pq" // PostgreSQL driver — the underscore means "import for side effects only"
)

func main() {
	// 1. Load configuration from environment
	cfg := config.Load()

	// 2. Connect to PostgreSQL
	conn, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Can't connect to the database:", err)
	}

	// 3. Create dependencies
	db := database.New(conn) // anathi hu badhi sql queries run kravi sakis.
	h := handler.New(db, cfg.JWTSecret)

	// 4. Start background scraper
	go scraper.Start(db, cfg.ScrapeConcurrency, cfg.ScrapeInterval)

	// 5. Set up router
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// 6. Register routes
	v1 := chi.NewRouter()

	// Public routes
	v1.Get("/healthz", h.HandleReadiness)
	v1.Get("/err", h.HandleError)
	v1.Post("/register", h.HandleRegister)
	v1.Post("/login", h.HandleLogin)
	v1.Get("/feeds", h.HandleGetFeeds)

	// Authenticated routes
	v1.Get("/users", middleware.WithAuth(db, cfg.JWTSecret, h.HandleGetUser))
	v1.Post("/feeds", middleware.WithAuth(db, cfg.JWTSecret, h.HandleCreateFeed))
	v1.Post("/feed_follows", middleware.WithAuth(db, cfg.JWTSecret, h.HandleCreateFeedFollow))
	v1.Get("/feed_follows", middleware.WithAuth(db, cfg.JWTSecret, h.HandleGetFeedFollows))
	v1.Delete("/feed_follows/{feedFollowID}", middleware.WithAuth(db, cfg.JWTSecret, h.HandleDeleteFeedFollow))
	v1.Get("/posts", middleware.WithAuth(db, cfg.JWTSecret, h.HandleGetPostsForUser))

	// Phase 4 changes
	v1.Put("/posts/{postID}/read", middleware.WithAuth(db, cfg.JWTSecret, h.HandleMarkPostAsRead))
	v1.Delete("/posts/{postID}/read", middleware.WithAuth(db, cfg.JWTSecret, h.HandleMarkPostAsUnread))
	// Bookmark routes
	v1.Put("/posts/{postID}/bookmark", middleware.WithAuth(db, cfg.JWTSecret, h.HandleBookmarkPost))
	v1.Delete("/posts/{postID}/bookmark", middleware.WithAuth(db, cfg.JWTSecret, h.HandleUnbookmarkPost))
	v1.Get("/bookmarks", middleware.WithAuth(db, cfg.JWTSecret, h.HandleGetBookmarks))

	router.Mount("/v1", v1)

	// 7. Start HTTP server
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + cfg.Port,
	}

	log.Printf("Server starting on port: %v", cfg.Port)
	log.Fatal(srv.ListenAndServe())
}
