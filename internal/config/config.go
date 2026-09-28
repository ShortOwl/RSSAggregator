package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application.
// Values are loaded from environment variables.
type Config struct {
	Port              string
	DatabaseURL       string
	JWTSecret         string
	SMTPEmail         string
	SMTPPassword      string
	ScrapeInterval    time.Duration
	ScrapeConcurrency int
	PostRetentionDays int32
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

	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	if smtpEmail == "" || smtpPassword == "" {
		log.Fatal("SMTP_EMAIL AND SMTP_PASSWORD are required")
	}
	retentionDays, err := parseRetentionDays(os.Getenv("POST_RETENTION_DAYS"))
	if err != nil {
		log.Fatal(err)
	}

	return Config{
		Port:              port,
		DatabaseURL:       dbURL,
		JWTSecret:         jwtSecret,
		SMTPEmail:         smtpEmail,
		SMTPPassword:      smtpPassword,
		ScrapeInterval:    10 * time.Minute,
		ScrapeConcurrency: 10,
		PostRetentionDays: retentionDays,
	}
}

func parseRetentionDays(value string) (int32, error) {
	if value == "" {
		return 50, nil // Limit how long articles occupy database storage by default.
	}
	days, err := strconv.ParseInt(value, 10, 32)
	if err != nil || days <= 0 {
		return 0, fmt.Errorf("POST_RETENTION_DAYS must be a positive 32-bit integer")
	}
	return int32(days), nil
}
