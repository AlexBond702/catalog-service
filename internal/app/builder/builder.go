package builder

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/AlexBond702/catalog-service/internal/app/config"
	rhandler "github.com/AlexBond702/catalog-service/internal/app/handler/http"
	hhealth "github.com/AlexBond702/catalog-service/internal/app/handler/http/health"
	"github.com/AlexBond702/catalog-service/internal/app/processor"
	pprocessor "github.com/AlexBond702/catalog-service/internal/app/processor/other"
	"github.com/AlexBond702/catalog-service/internal/app/repository"
	pcategory "github.com/AlexBond702/catalog-service/internal/app/repository/category"
	rcpostgres "github.com/AlexBond702/catalog-service/internal/app/repository/conn/postgres"
	pproduct "github.com/AlexBond702/catalog-service/internal/app/repository/product"
)

type Builder struct {
	cCtx *cli.Context
	ctx  context.Context
	wg   sync.WaitGroup
	err  error
	cfg  config.Config

	connPostgres *rcpostgres.Client

	categoryRepo repository.Category
	productRepo  repository.Product

	healthHandler rhandler.Health
	// categoryHandler rhandler.Category
	// productHandler  rhandler.Product

	processors []processor.Processor
}

func NewBuilder(cCtx *cli.Context) *Builder {
	b := Builder{
		cCtx: cCtx,
		ctx:  context.Background(),
	}
	b.healthHandler = hhealth.NewHandler()
	return &b
}

func (b *Builder) BuildConfig(injectors ...func(c *config.Config)) {
	b.exec(true, func(b *Builder) {
		b.buildConfig(config.LoadArgs{}, injectors)
	})
}

func (b *Builder) BuildConfigSimple(injectors ...func(c *config.Config)) {
	b.exec(true, func(b *Builder) {
		b.buildConfig(config.LoadArgs{SkipConfig: true}, injectors)
	})
}

func (b *Builder) Run() {
	if b.err != nil {
		log.Fatal().Err(b.err).Msg("Failed to initialize application")
	} else {
		log.Info().Msg("Application is initialized")
	}
	defer log.Info().Msg("Application is completed, GoodBye!")

	for _, proc := range b.processors {
		proc.StartAsync(b.ctx, &b.wg)
	}
	b.wg.Wait()
}

func (b *Builder) BuildRepoConnPostgres() {
	b.exec(true, func(b *Builder) {
		configDB := b.cfg.Repository.Postgres
		client, err := rcpostgres.NewConn(b.ctx, configDB)
		if err != nil {
			b.err = err
			return
		}
		b.connPostgres = client
	})
}

func (b *Builder) BuildRepoConnMigrate() {
	b.exec(b.connPostgres != nil, func(b *Builder) {
		proc := pprocessor.NewMigrator(b.connPostgres)
		b.processors = append(b.processors, proc)
	})
}

func (b *Builder) BuildRepoCategory() {
	b.exec(true, func(b *Builder) {
		repoCategory := pcategory.NewRepoFromPostgres(b.ctx, b.connPostgres)
		b.categoryRepo = repoCategory
	}, b.connPostgres)
}

func (b *Builder) BuildRepoProduct() {
	b.exec(true, func(b *Builder) {
		repoProduct := pproduct.NewRepoFromPostgres(b.ctx, b.connPostgres)
		b.productRepo = repoProduct
	}, b.connPostgres)
}

func (b *Builder) buildConfig(args config.LoadArgs, injectors []func(*config.Config)) {
	args.Output = b.cCtx.App.Writer
	args.EnableSimpleLog = b.cCtx.Bool("no json")

	config.Load(args)

	for i, injector := range injectors {
		if injectors[i] != nil {
			injector(&config.Root)
		}
	}
	b.cfg = config.Root
}

func (b *Builder) exec(preCond bool, cb func(b *Builder), requiredArgs ...any) {
	if !preCond || b.err != nil {
		return
	}

	for i, requiredArg := range requiredArgs {
		rv := reflect.ValueOf(requiredArg)
		if !rv.IsValid() {
			b.err = fmt.Errorf("BUG: required argument #%d is nil (check dependencies)", i)
			return
		}
		if rv.Type().Kind() == reflect.Struct || !rv.IsZero() {
			continue
		}
		b.err = fmt.Errorf("BUG: required %s, but empty", rv.Type().String())
		return
	}
	cb(b)
}
