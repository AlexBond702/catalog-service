package main

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/AlexBond702/catalog-service/internal/app/config"
	rhealth "github.com/AlexBond702/catalog-service/internal/app/handler/health"
	rprocessor "github.com/AlexBond702/catalog-service/internal/app/processor/http"
	rcpostgres "github.com/AlexBond702/catalog-service/internal/app/repository/conn/postgres"
)

func main() {
	config.Load()
	cfg := config.Root
	log.Print("migration table:", cfg.Repository.Postgres.MigrationTable)
	ctx := context.Background()

	pgClient, err := rcpostgres.NewConn(ctx, cfg.Repository.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}

	// Применение миграций
	oldVer, newVer, err := pgClient.Migrate(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}

	if oldVer != newVer {
		log.Info().
			Int64("old_version", oldVer).
			Int64("new_version", newVer).
			Msg("Database migrated")
	} else {
		log.Info().
			Int64("version", newVer).
			Msg("Database is up to date")
	}

	hHealth := rhealth.NewHandler()

	httpServer := rprocessor.NewHttp(hHealth, cfg.Processor.WebServer)
	if err := httpServer.Serve(); err != nil {
		log.Fatal().Err(err).Msg("HTTP server failed")
	}
}
