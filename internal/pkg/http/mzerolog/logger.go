package mzerolog

import (
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/AlexBond702/catalog-service/internal/pkg/http/httph"
)

type middleware struct {
	log zerolog.Logger

	fromOption struct {
		skipper func(r *http.Request) bool
	}
}

func (m *middleware) Callback(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const (
			TailSuccess = " finished with no error"
			TailFail    = " finished (or aborted) with error"
		)
		start := time.Now()
		next.ServeHTTP(w, r)
		err := httph.ErrorGet(r)
		execTime := time.Since(start)

		if m.fromOption.skipper(r) {
			return
		}

		var mb strings.Builder
		mb.Grow(48 + len(r.RequestURI))
		mb.WriteString(r.Method)
		mb.WriteByte(' ')
		mb.WriteString(r.RequestURI)

		var ev *zerolog.Event

		if err == nil {
			ev = m.log.Debug()
			mb.WriteString(TailSuccess)
		} else {
			ev = m.log.Error()
			mb.WriteString(TailFail)
		}
		ev.Err(err).
			Str("exec_time", execTime.String()).
			Str("client_ip", r.RemoteAddr).
			Msg(mb.String())
	})
}
