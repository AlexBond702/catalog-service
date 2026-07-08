package grpc

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	pb "github.com/AlexBond702/catalog-service/gen/proto/catalog/v1"
	"github.com/AlexBond702/catalog-service/internal/app/config/section"
	"github.com/AlexBond702/catalog-service/internal/app/handler/grpc/catalog"
	"github.com/AlexBond702/catalog-service/internal/app/processor"
	"github.com/AlexBond702/catalog-service/internal/app/util"
	"github.com/AlexBond702/catalog-service/internal/pkg/grpc/grpch"
	"github.com/AlexBond702/catalog-service/internal/pkg/grpc/mcommon"
	"github.com/AlexBond702/catalog-service/internal/pkg/grpc/mprom"
	"github.com/AlexBond702/catalog-service/internal/pkg/grpc/mzerolog"
)

type grpcProc struct {
	server *grpc.Server
	addr   string
}

func NewGrpc(handler catalog.Handler,
	unaryInterceptors []grpc.UnaryServerInterceptor,
	streamInterceptors []grpc.StreamServerInterceptor,
	cfg section.ProcessorGrpc,
) processor.Processor {
	coreMiddlewares := []grpch.Middleware{
		mcommon.NewRecovery(),
		mzerolog.NewMiddleware(),
		mprom.New(),
	}

	coreUnary := make([]grpc.UnaryServerInterceptor, 0)
	coreStream := make([]grpc.StreamServerInterceptor, 0)

	for _, m := range coreMiddlewares {
		coreUnary = append(coreUnary, m.ForUnary())
		coreStream = append(coreStream, m.ForStream())
	}

	unaryInterceptors = append(coreUnary, unaryInterceptors...)
	streamInterceptors = append(coreStream, streamInterceptors...)

	mprom.EnableHistogram()

	server := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)
	pb.RegisterCatalogServiceServer(server, &handler)

	return &grpcProc{
		server: server,
		addr:   fmt.Sprintf("%s:%d", cfg.Host, cfg.ListenPort),
	}
}

func (g *grpcProc) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	var lc net.ListenConfig
	l, err := lc.Listen(ctx, "tcp", g.addr)
	if err != nil {
		log.Error().Err(err).Str("listen_addr:", g.addr).
			Msg("Failed to start listening TCP addr for gRPC server")
		return
	}
	log.Info().Str("listen_addr", g.addr).
		Msg("Listening of TCP addr for gRPC server has been started")
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := g.server.Serve(l); err != nil {
			log.Error().Err(err).Msg("failed to Serve gRPC")
			return
		}
	}()
	go func() {
		processor.WatchForShutdown(ctx, wg, util.CloserFunc(l.Close))
	}()
	go func() {
		_ = util.CloserFunc(func() error {
			g.server.GracefulStop()
			return nil
		})
	}()
}
