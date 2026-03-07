package main

import (
	"github.com/rs/zerolog/log"

	"github.com/AlexBond702/catalog-service/internal/app/config"
	rhealth "github.com/AlexBond702/catalog-service/internal/app/handler/health"
	rprocessor "github.com/AlexBond702/catalog-service/internal/app/processor/http"
)

func main() {
	config.Load()

	cfg := config.Root

	hHealth := rhealth.NewHandler()

	httpServer := rprocessor.NewHttp(hHealth, cfg.Processor.WebServer)
	if err := httpServer.Serve(); err != nil {
		log.Fatal().Err(err).Msg("HTTP server failed")
	}
}
