package builder

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"

	"github.com/AlexBond702/catalog-service/internal/app/config"
	"github.com/AlexBond702/catalog-service/internal/app/handler/grpc/catalog"
	rhandler "github.com/AlexBond702/catalog-service/internal/app/handler/http"
	hcategory "github.com/AlexBond702/catalog-service/internal/app/handler/http/category"
	hhealth "github.com/AlexBond702/catalog-service/internal/app/handler/http/health"
	hproduct "github.com/AlexBond702/catalog-service/internal/app/handler/http/product"
	"github.com/AlexBond702/catalog-service/internal/app/monitor/metric"
	"github.com/AlexBond702/catalog-service/internal/app/processor"
	"github.com/AlexBond702/catalog-service/internal/app/processor/gateway"
	pgrpc "github.com/AlexBond702/catalog-service/internal/app/processor/grpc"
	rprocessor "github.com/AlexBond702/catalog-service/internal/app/processor/http"
	pprocessor "github.com/AlexBond702/catalog-service/internal/app/processor/other"
	"github.com/AlexBond702/catalog-service/internal/app/repository"
	pcategory "github.com/AlexBond702/catalog-service/internal/app/repository/category"
	rcpostgres "github.com/AlexBond702/catalog-service/internal/app/repository/conn/postgres"
	pproduct "github.com/AlexBond702/catalog-service/internal/app/repository/product"
	"github.com/AlexBond702/catalog-service/internal/app/service"
	scategory "github.com/AlexBond702/catalog-service/internal/app/service/category"
	sproduct "github.com/AlexBond702/catalog-service/internal/app/service/product"
)

type Builder struct {
	cCtx    *cli.Context
	ctx     context.Context
	wg      sync.WaitGroup
	err     error
	cfg     config.Config
	chError chan error

	connPostgres *rcpostgres.Client

	categoryRepo repository.Category
	productRepo  repository.Product

	categoryService service.Category
	productService  service.Product

	healthHandler   rhandler.Health
	categoryHandler rhandler.Category
	productHandler  rhandler.Product

	interceptorsUnary  []grpc.UnaryServerInterceptor
	interceptorsStream []grpc.StreamServerInterceptor

	grpcCatalogHandler *catalog.Handler

	processors []processor.Processor
}

func NewBuilder(cCtx *cli.Context) *Builder {
	b := Builder{
		cCtx:    cCtx,
		ctx:     context.Background(),
		chError: make(chan error, 4096),
	}

	ctxCancel, cancelFunc := context.WithCancel(context.Background())
	b.ctx = ctxCancel

	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go b.waitForSignal(chanSignal, cancelFunc)
	go b.printErrors()

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

func (b *Builder) BuildServiceCategory() {
	b.exec(true, func(b *Builder) {
		serviceCat := scategory.NewService(b.categoryRepo, b.productRepo)
		b.categoryService = serviceCat
	}, b.categoryRepo)
}

func (b *Builder) BuildServiceProduct() {
	b.exec(true, func(b *Builder) {
		serviceProd := sproduct.NewService(b.productRepo, b.categoryRepo)
		b.productService = serviceProd
	}, b.productRepo)
}

func (b *Builder) BuildHandlerHttpCategory() {
	b.exec(true, func(b *Builder) {
		categoryHandler := hcategory.NewHandler(b.categoryService)
		b.categoryHandler = categoryHandler
	}, b.categoryService)
}

func (b *Builder) BuildHandlerHttpProduct() {
	b.exec(true, func(b *Builder) {
		productHandler := hproduct.NewHandler(b.productService)
		b.productHandler = productHandler
	}, b.productService)
}

func (b *Builder) BuildHandlerGrpcCatalog() {
	b.exec(true, func(b *Builder) {
		b.grpcCatalogHandler = catalog.NewHandler(b.productService)
	}, b.productService)
}

func (b *Builder) BuildProcHttp() {
	b.exec(true, func(b *Builder) {
		procHttp := rprocessor.NewHttp(b.healthHandler, b.categoryHandler, b.productHandler, nil, b.cfg.Processor.WebServer)
		b.processors = append(b.processors, procHttp)
	}, b.productHandler, b.categoryHandler)
}

func (b *Builder) BuildProcGrpc() {
	b.exec(true, func(b *Builder) {
		procGrpc := pgrpc.NewGrpc(*b.grpcCatalogHandler, b.interceptorsUnary, b.interceptorsStream, b.cfg.Processor.Grpc)
		b.processors = append(b.processors, procGrpc)
	}, b.grpcCatalogHandler)
}

func (b *Builder) BuildProcGrpcGateway() {
	b.exec(true, func(b *Builder) {
		grpcGateway := gateway.NewGateway(b.cfg.Processor.Gateway)
		b.processors = append(b.processors, grpcGateway)
	})
}

func (b *Builder) BuildMonitorPrometheus() {
	b.exec(true, func(b *Builder) {
		if !b.cfg.Monitor.Prometheus.Enabled {
			log.Warn().Msg("Prometheus metrics disabled")
			return
		}
		prometheus := metric.NewPrometheusObserver()
		b.processors = append(b.processors, prometheus)
	})
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

func (b *Builder) waitForSignal(sig chan os.Signal, cancelFunc func()) {
	sigValue := <-sig
	log.Printf("Signal arrive: %T", sigValue)
	cancelFunc()
}

func (b *Builder) printErrors() {
	for chErr := range b.chError {
		log.Error().Err(chErr)
	}
}
