package pprocessor

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/AlexBond702/catalog-service/internal/app/processor"
	"github.com/AlexBond702/catalog-service/internal/app/repository"
)

type procMigrate struct {
	migrate repository.Migrate
}

func NewMigrator(migrator repository.Migrate) processor.Processor {
	return &procMigrate{migrator}
}

func (p *procMigrate) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	processor.Wrap(ctx, wg, p.job)
}

func (p *procMigrate) job(ctx context.Context) {
	oldVer, newVer, err := p.migrate.Migrate(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to migrate repository schema")
		return
	}
	if oldVer != newVer {
		log.Info().Int64("old_ver", oldVer).Int64("new_ver", newVer).Msg("Repository schema has been updated")
	} else {
		log.Info().Msg("Repository schema is up to date, nothing to migrate")
	}
}
