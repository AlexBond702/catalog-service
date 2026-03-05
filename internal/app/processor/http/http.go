package rprocessor

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/AlexBond702/catalog-service/internal/app/config/section"
	rhandler "github.com/AlexBond702/catalog-service/internal/app/handler"
)

type httpProc struct {
	server http.Server
	addr   string
}

func NewHttp(hHealth rhandler.Health, cfg section.ProcessorWebServer) *httpProc {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(handlerNotFound)
	vGenericRegHealthCheck(r, hHealth)

	p := httpProc{
		addr: fmt.Sprintf(":%d", cfg.ListenPort),
	}
	_ = r.Walk(func(route *mux.Route, router *mux.Router, sl []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		if path == "" {
			log.Debug().Msg("Empty path")
			return nil
		}
		log.Debug().Str("Path", path).Msg("Path")
		return nil
	})
	_ = r.Walk(func(route *mux.Route, router *mux.Router, sl []*mux.Route) error {
		methods, _ := route.GetMethods()
		if len(methods) == 0 {
			log.Debug().Msg("Empty method")
			return nil
		}
		log.Debug().Strs("Methods", methods).Msg("Methods")
		return nil
	})
	p.server.Addr = p.addr
	p.server.Handler = r
	return &p
}

func (p *httpProc) Serve() error {
	log.Info().Str("addr", p.addr).Msg("Starting HTTP server")
	return p.server.ListenAndServe()
}
