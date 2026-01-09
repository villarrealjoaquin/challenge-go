package config

import (
	"os"
)

type Config struct {
	BooksAPIURL string
	Port        string
}

func Load() *Config {
	return &Config{
		BooksAPIURL: getEnv("BOOKS_API_URL", "https://6781684b85151f714b0aa5db.mockapi.io/api/v1/books"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
