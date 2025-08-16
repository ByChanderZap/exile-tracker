package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ByChanderZap/exile-tracker/cmd/api"
	"github.com/ByChanderZap/exile-tracker/config"
	"github.com/ByChanderZap/exile-tracker/db"
	"github.com/ByChanderZap/exile-tracker/utils"
	"github.com/rs/zerolog"
)

func main() {
	log := utils.ChildLogger("main")

	server := api.NewAPIServer(config.Envs.Port)

	db, err := db.NewSqliteStorage(config.Envs.DBPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}

	initStorage(db, log)

	// Start server in a goroutine
	go func() {
		if err := server.Start(); err != nil && err.Error() != "http: Server closed" {
			log.Fatal().Err(err).Msg("Failed to start API server")
		}
	}()

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Info().Msg("Shutting down application...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop API server gracefully
	if err := server.Stop(ctx); err != nil {
		log.Error().Err(err).Msg("Error shutting down API server")
	}

	log.Info().Msg("Application shutdown complete.")
}

func initStorage(db *sql.DB, log zerolog.Logger) {
	err := db.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("Database ping failed")
	}
	log.Info().Msg("Database connection established")

}
