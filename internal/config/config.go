package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr string
}

func Load() *Config {
	_ = godotenv.Load()

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	return &Config{
		HTTPAddr: addr,
	}
}
