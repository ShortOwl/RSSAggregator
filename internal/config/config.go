package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application.
// Values are loaded from environment variables.
type Config struct {
	Port              string
	DatabaseURL       string
	JWTSecret         string
	ScrapeInterval    time.Duration
	ScrapeConcurrency int
}

// Load reads configuration from environment variables (and .env file).
// It terminates the program if required variables are missing.
func Load() Config {
	// Load .env file if it exists (ignored in production where
	// env vars are set directly by the deployment platform).
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not found in the environment")
	}
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters")
	}

	return Config{
		Port:              port,
		DatabaseURL:       dbURL,
		JWTSecret:         jwtSecret,
		ScrapeInterval:    10 * time.Minute,
		ScrapeConcurrency: 10,
	}
}
