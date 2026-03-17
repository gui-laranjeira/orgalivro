package config

import (
	"os"
	"strings"
)

type Config struct {
	DBPath         string
	Port           string
	AllowedOrigins []string
}

func Load() Config {
	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins == "" {
		origins = "http://localhost:5173"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./orgalivro.db"
	}
	return Config{
		DBPath:         dbPath,
		Port:           port,
		AllowedOrigins: strings.Split(origins, ","),
	}
}
