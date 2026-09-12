package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/ShortOwl/RSSAggregator/internal/config"
	"github.com/ShortOwl/RSSAggregator/internal/database"
	"github.com/ShortOwl/RSSAggregator/internal/handler"
	"github.com/ShortOwl/RSSAggregator/internal/middleware"
	"github.com/ShortOwl/RSSAggregator/internal/scraper"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"golang.org/x/time/rate"

	_ "github.com/lib/pq" // PostgreSQL driver — the underscore means "import for side effects only"
)

func main() {
	// 1. Load configuration from environment
	cfg := config.Load()

	var logHandler slog.Handler = slog.NewTextHandler(os.Stdout, nil)
	if os.Getenv("ENVIRONMENT") == "production" {
		logHandler = slog.NewJSONHandler(os.Stdout, nil)
	}
	logger := slog.New(logHandler) // Handler decides output format like JSON,txt. And Logger recieves the events to write.
	slog.SetDefault(logger)

	// 2. Connect to PostgreSQL
	conn, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Error("can't connect to the database", "error", err)
		// logger.Error(message, key, value)
		os.Exit(1)
	}

	// 3. Create dependencies
	db := database.New(conn) // anathi hu badhi sql queries run kravi sakis.
	h := handler.New(db, cfg.JWTSecret)

	// 4. Start background scraper
	go scraper.Start(db, cfg.ScrapeConcurrency, cfg.ScrapeInterval)

	// 5. Set up router
	router := chi.NewRouter()
	rateLimiter := middleware.NewRateLimiter(rate.Limit(5), 10)
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recovery)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link", "X-Request-ID"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	router.Use(rateLimiter.Limit)

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

	logger.Info("Server starting", "port", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Server stopped", "error", err)
		os.Exit(1)
	}
}
