package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Promise111/go-rest-api-todo/internal/config"
	"github.com/Promise111/go-rest-api-todo/internal/database"
	"github.com/Promise111/go-rest-api-todo/internal/router"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	fmt.Println("🚀 Todo API Server started!")
	fmt.Println("💹 Project setup complete")
	fmt.Println("📁 Folder structure all set")

	var cfg *config.Config
	var configError error

	cfg, configError = config.Load()

	if configError != nil {
		slog.Error("Unable to load configuration")
		os.Exit(1)
	}

	var pool *pgxpool.Pool
	pool, err := database.Connect(cfg.DatabaseURL)

	if err != nil {
		slog.Error("Failed to connect to Database")
		defer pool.Close()
		os.Exit(1)
	}

	r := router.Router()

	r.Run(":" + cfg.Port)
}
