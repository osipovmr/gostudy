package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr string
	DBURL    string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		HTTPAddr: getEnv("HTTP_ADDR", ":8080"),
		DBURL:    getEnv("DB_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
