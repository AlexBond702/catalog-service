package gateway

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/AlexBond702/catalog-service/gen/proto/catalog/v1"
	"github.com/AlexBond702/catalog-service/internal/app/config/section"
	"github.com/AlexBond702/catalog-service/internal/app/processor"
	"github.com/AlexBond702/catalog-service/internal/app/util"
)

type gateWay struct {
	server       *http.Server
	addr         string
	grpcEndPoint string
}

func NewGateway(cfg section.ProcessorGateway) processor.Processor {
	addr := fmt.Sprintf("%s:%v", cfg.Host, cfg.ListenPort)
	return &gateWay{
		addr:         addr,
		grpcEndPoint: cfg.GrpcEndpoint,
	}
}

func (g *gateWay) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	gwMux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	err := pb.RegisterCatalogServiceHandlerFromEndpoint(ctx, gwMux, g.grpcEndPoint, opts)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to register gRPC-Gateway")
		return
	}
	g.server = &http.Server{
		Addr:              g.addr,
		Handler:           gwMux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	var lc net.ListenConfig
	l, err := lc.Listen(ctx, "tcp", g.addr)
	if err != nil {
		log.Error().Err(err).Str("listen_addr", g.addr).
			Msg("Failed to start listening TCP addr for gRPC-Gateway")
		return
	}
	log.Info().Str("listen_addr", g.addr).
		Msg("Listening of TCP addr for gRPC-Gateway has been started")
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := g.server.Serve(l); err != nil {
			log.Error().Err(err).Msg("failed to Serve GateWay")
		}
	}()
	go processor.WatchForShutdown(ctx, wg, util.CloserFunc(l.Close))
	go processor.WatchForShutdown(ctx, wg, util.NewCloserContextFunc(g.server.Shutdown, context.Background(), 5*time.Second))
}
