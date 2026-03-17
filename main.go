package main

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/AlexBond702/catalog-service/internal/app/config"
	hhealth "github.com/AlexBond702/catalog-service/internal/app/handler/health"
	hcategory "github.com/AlexBond702/catalog-service/internal/app/handler/http/category"
	"github.com/AlexBond702/catalog-service/internal/app/handler/http/product"
	rprocessor "github.com/AlexBond702/catalog-service/internal/app/processor/http"
	pcategory "github.com/AlexBond702/catalog-service/internal/app/repository/category"
	rcpostgres "github.com/AlexBond702/catalog-service/internal/app/repository/conn/postgres"
	pproduct "github.com/AlexBond702/catalog-service/internal/app/repository/product"
	"github.com/AlexBond702/catalog-service/internal/app/service/category"
	sproduct "github.com/AlexBond702/catalog-service/internal/app/service/product"
)

func main() {
	config.Load()
	cfg := config.Root
	ctx := context.Background()

	pgClient, err := rcpostgres.NewConn(ctx, cfg.Repository.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to PostgresSQL")
	}

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
	// репо
	categoryRepo := pcategory.NewRepoFromPostgres(pgClient)
	productRepo := pproduct.NewRepoFromPostgres(pgClient)

	// сервис
	categorySVC := scategory.NewService(categoryRepo, productRepo)
	productSVC := sproduct.NewService(productRepo, categoryRepo)

	// хендлеры
	healthHandler := hhealth.NewHandler()
	categoryHandler := hcategory.NewHandler(categorySVC)
	productHandler := hproduct.NewHandler(productSVC)

	httpServer := rprocessor.NewHttp(healthHandler, cfg.Processor.WebServer, categoryHandler, productHandler)
	if err := httpServer.Serve(); err != nil {
		log.Fatal().Err(err).Msg("HTTP server failed")
	}
}
