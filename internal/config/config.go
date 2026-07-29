package config

import (
	"errors"
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
		slog.Info("Fetch env variables", "Warning: ", ".env file not found, using environment variables")
		return nil, errors.New(".env file not found.")
	}

	var config *Config = &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
	}

	return config, nil
}
