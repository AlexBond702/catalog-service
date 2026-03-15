package main

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/AlexBond702/catalog-service/internal/app/config"
	"github.com/AlexBond702/catalog-service/internal/app/entity"
	rhealth "github.com/AlexBond702/catalog-service/internal/app/handler/health"
	rprocessor "github.com/AlexBond702/catalog-service/internal/app/processor/http"
	pcategory "github.com/AlexBond702/catalog-service/internal/app/repository/category"
	rcpostgres "github.com/AlexBond702/catalog-service/internal/app/repository/conn/postgres"
	pproduct "github.com/AlexBond702/catalog-service/internal/app/repository/product"
)

func main() {
	config.Load()
	cfg := config.Root
	ctx := context.Background()

	pgClient, err := rcpostgres.NewConn(ctx, cfg.Repository.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
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

	categoryRepo := pcategory.NewRepoFromPostgres(pgClient)
	productRepo := pproduct.NewRepoFromPostgres(pgClient)

	cat := entity.Category{
		GUID:      uuid.Must(uuid.NewRandom()),
		Name:      "Электроника",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := categoryRepo.Create(ctx, cat); err != nil {
		log.Fatal().Err(err).Msg("Category create failed")
	}
	log.Info().Str("guid", cat.GUID.String()).Msg("Category created")

	found, err := categoryRepo.GetByGUID(ctx, cat.GUID)
	if err != nil {
		log.Fatal().Err(err).Msg("Category GetByGUID failed")
	}
	log.Info().Str("name", found.Name).Msg("Category found")

	found.Name = "Бытовая техника"
	found.UpdatedAt = time.Now()
	if err := categoryRepo.Update(ctx, found); err != nil {
		log.Fatal().Err(err).Msg("Category update failed")
	}
	log.Info().Msg("Category updated")

	allCats, err := categoryRepo.List(ctx, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Category list failed")
	}
	log.Info().Int("count", len(allCats)).Msg("All categories")

	filterName := "Бытовая техника"
	filtered, err := categoryRepo.List(ctx, &filterName)
	if err != nil {
		log.Fatal().Err(err).Msg("Category list by name failed")
	}
	log.Info().Int("count", len(filtered)).Msg("Filtered categories")

	desc := "Мощный пылесос"
	prod := entity.Product{
		GUID:         uuid.Must(uuid.NewRandom()),
		Name:         "Пылесос Dyson V15",
		Description:  &desc,
		Price:        49999.99,
		CategoryGUID: cat.GUID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := productRepo.Create(ctx, prod); err != nil {
		log.Fatal().Err(err).Msg("Product create failed")
	}
	log.Info().Str("guid", prod.GUID.String()).Msg("Product created")

	products, err := productRepo.List(ctx, nil, &cat.GUID)
	if err != nil {
		log.Fatal().Err(err).Msg("Product list by category failed")
	}
	log.Info().Int("count", len(products)).Msg("Products in category")

	if err := productRepo.Delete(ctx, prod.GUID); err != nil {
		log.Fatal().Err(err).Msg("Product delete failed")
	}
	if err := categoryRepo.Delete(ctx, cat.GUID); err != nil {
		log.Fatal().Err(err).Msg("Category delete failed")
	}
	log.Info().Msg("Cleanup complete")

	hHealth := rhealth.NewHandler()

	httpServer := rprocessor.NewHttp(hHealth, cfg.Processor.WebServer)
	if err := httpServer.Serve(); err != nil {
		log.Fatal().Err(err).Msg("HTTP server failed")
	}
}
