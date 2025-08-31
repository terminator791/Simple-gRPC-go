package config

import (
	"os"
)

type Config struct {
	DatabaseURL     string
	Port           string
	UserServiceAddr string
}

func Load() *Config {
	return &Config{
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5433/order_service?sslmode=disable"),
		Port:           getEnv("PORT", "50052"),
		UserServiceAddr: getEnv("USER_SERVICE_ADDR", "localhost:50051"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}