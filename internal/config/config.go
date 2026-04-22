package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr      string
	DBURL         string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	IdleTimeout   time.Duration
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		HTTPAddr:      getEnv("HTTP_ADDR", ":8080"),
		DBURL:         getEnv("DB_URL", ""),
		ReadTimeout:   mustDuration(getEnv("READ_TIMEOUT", "5s")),
		WriteTimeout:  mustDuration(getEnv("WRITE_TIMEOUT", "10s")),
		AccessSecret:  getEnv("AccessSecret", "AccessSecret"),
		RefreshSecret: getEnv("RefreshSecret", "RefreshSecret"),
		AccessTTL:     mustDuration(getEnv("AccessTTL", "3000s")),
		RefreshTTL:    mustDuration(getEnv("RefreshTTL", "300000s")),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustDuration(value string) time.Duration {
	d, err := time.ParseDuration(value)
	if err != nil {
		panic(err)
	}
	return d
}
