package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for our application
type Config struct {
	Port         string
	DatabaseURL  string
	AllowOrigins string
}

// Load returns configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	// Set default values
	config := &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		AllowOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:5173"),
	}

	// Validate required DATABASE_URL
	if config.DatabaseURL == "" {
		panic("DATABASE_URL environment variable is required")
	}

	return config, nil
}

// Helper function to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
