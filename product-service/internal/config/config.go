package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

// Load loads configuration from environment variables
func Load() *Config {
	config := &Config{
		Port:        getEnv("PRODUCT_SERVICE_PORT", "50053"),
		DatabaseURL: getEnv("PRODUCT_DATABASE_URL", "postgres://product_user:product_pass@localhost:5434/product_db?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
	}
	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}