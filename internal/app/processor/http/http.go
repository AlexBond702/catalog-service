package rprocessor

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/AlexBond702/catalog-service/internal/app/config/section"
	rhandler "github.com/AlexBond702/catalog-service/internal/app/handler/http"
)

type httpProc struct {
	server http.Server
	addr   string
}

func NewHttp(hHealth rhandler.Health, cfg section.ProcessorWebServer, hCategory rhandler.Category, hProduct rhandler.Product) *httpProc {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(handlerNotFound)
	vGenericRegHealthCheck(r, hHealth)
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

func (p *httpProc) Serve() error {
	log.Info().Str("addr", p.addr).Msg("Starting HTTP server")
	return p.server.ListenAndServe()
}
