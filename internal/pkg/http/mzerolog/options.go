package mzerolog

import (
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/AlexBond702/catalog-service/internal/pkg/http/httph"
)

type Option = func(m *middleware)

func WithLogger(l zerolog.Logger) Option {
	return func(m *middleware) {
		m.log = l
	}
}

func WithSkipper(skipper func(r *http.Request) bool) Option {
	if skipper != nil {
		return func(m *middleware) {
			m.fromOption.skipper = skipper
		}
	}
	return nil
}

func NewMiddleware(opts ...Option) httph.Middleware {
	m := middleware{
		log: log.Logger,
	}
	m.fromOption.skipper = defaultSkipper
	for _, opt := range opts {
		opt(&m)
	}
	return m.Callback
}

func defaultSkipper(_ *http.Request) bool {
	return false
}
