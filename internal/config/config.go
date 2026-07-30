package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func Load() (*Config, error) {
	var err error = godotenv.Load()

	if err != nil {
		slog.Error("Warning: ", "message", ".env file not found "+err.Error())
		return nil, err
	}

	var config *Config = &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
	}

	return config, nil
}
