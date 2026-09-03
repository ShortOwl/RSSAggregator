package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/ShortOwl/RSSAggregator/internal/config"
	"github.com/ShortOwl/RSSAggregator/internal/database"
	"github.com/ShortOwl/RSSAggregator/internal/handler"
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
	h := handler.New(db)

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
	v1.Post("/users", h.HandleCreateUser)
	v1.Get("/feeds", h.HandleGetFeeds)

	// Authenticated routes
	v1.Get("/users", h.WithAuth(h.HandleGetUser))
	v1.Post("/feeds", h.WithAuth(h.HandleCreateFeed))
	v1.Post("/feed_follows", h.WithAuth(h.HandleCreateFeedFollow))
	v1.Get("/feed_follows", h.WithAuth(h.HandleGetFeedFollows))
	v1.Delete("/feed_follows/{feedFollowID}", h.WithAuth(h.HandleDeleteFeedFollow))
	v1.Get("/posts", h.WithAuth(h.HandleGetPostsForUser))

	router.Mount("/v1", v1)

	// 7. Start HTTP server
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + cfg.Port,
	}

	log.Printf("Server starting on port: %v", cfg.Port)
	log.Fatal(srv.ListenAndServe())
}
