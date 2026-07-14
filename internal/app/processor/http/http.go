package rprocessor

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/AlexBond702/catalog-service/internal/app/config/section"
	rhandler "github.com/AlexBond702/catalog-service/internal/app/handler/http"
	"github.com/AlexBond702/catalog-service/internal/app/processor"
	"github.com/AlexBond702/catalog-service/internal/app/util"
	"github.com/AlexBond702/catalog-service/internal/pkg/http/httph"
	"github.com/AlexBond702/catalog-service/internal/pkg/http/mzerolog"
)

type httpProc struct {
	processor.Processor
	server http.Server
	addr   string
}

func extractUserID(r *http.Request) string {
	return r.Header.Get("X-User-ID")
}

func extractRequestID(r *http.Request) string {
	return r.Header.Get("X-Request-ID")
}

func extractContentLength(r *http.Request) any {
	if r.ContentLength > 0 {
		return r.ContentLength
	}
	return nil
}

func extractQuery(r *http.Request) any {
	return r.URL.RawQuery
}

func NewHttp(hHealth rhandler.Health,
	hCategory rhandler.Category,
	hProduct rhandler.Product,
	middlewares []httph.Middleware,
	cfg section.ProcessorWebServer,
) processor.Processor {
	r := mux.NewRouter()
	r.StrictSlash(true)

	logMW := mzerolog.NewMiddleware(
		mzerolog.WithSkipper(util.IsFilteredHttpRoute),
		mzerolog.WithStringExtractor("user_id", extractUserID),
		mzerolog.WithStringExtractorOnFail("request_id", extractRequestID),
		mzerolog.WithAnyExtractorOnFail("query", extractQuery),
		mzerolog.WithAnyExtractorOnSuccess("content_length", extractContentLength),
	)
	r.Use(middlewaresToGorilla(middlewares)...)

	r.Use(
		httph.NewErrorMiddleware(),
		logMW,
		makeErrorMiddleware(),
	)

	r.NotFoundHandler = http.HandlerFunc(handlerNotFound)
	vGenericRegHealthCheck(r, hHealth)
	vGenericRegPprof(r)
	vGenericRegMetrics(r)
	rV1 := r.PathPrefix("/v1").Subrouter()
	if hCategory != nil {
		v1CategoryHandler(rV1, hCategory)
	}
	if hProduct != nil {
		v1ProductHandler(rV1, hProduct)
	}
	_ = r.Walk(func(route *mux.Route, router *mux.Router, sl []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		if len(methods) == 0 {
			log.Debug().Msg("Empty method")
			return nil
		}

		if path == "" {
			log.Debug().Msg("Empty path")
			return nil
		}
		log.Debug().Str("Path", path).Msg("Path")
		log.Debug().Strs("Methods", methods).Msg("Methods")
		return nil
	})

	p := httpProc{
		addr: fmt.Sprintf(":%d", cfg.ListenPort),
	}

	p.server.Addr = p.addr
	p.server.Handler = r
	return &p
}

func (p *httpProc) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	var lc net.ListenConfig
	l, err := lc.Listen(ctx, "tcp", p.addr)
	if err != nil {
		log.Fatal().Err(err).Str("listen_addr", p.addr).
			Msg("Failed to start listening TCP addr for HTTP server")
		return
	}
	log.Info().Str("listen_addr", p.addr).
		Msg("Listening of TCP addr for HTTP server has been started")
	wg.Add(1)
	go func() {
		defer wg.Done()
		p.serve(l)
	}()
	go processor.WatchForShutdown(ctx, wg, util.CloserFunc(l.Close))
	go processor.WatchForShutdown(ctx, wg, util.NewCloserContextFunc(p.server.Shutdown, context.Background(), time.Second*5))
}

func (p *httpProc) serve(l net.Listener) {
	log.Info().Str("addr", p.addr).Msg("Starting HTTP server")
	_ = p.server.Serve(l) // блокирует горутину
}
