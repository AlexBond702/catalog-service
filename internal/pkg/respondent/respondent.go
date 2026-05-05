package respondent

import (
	"errors"
	"net/http"

	"github.com/AlexBond702/catalog-service/internal/pkg/http/httph"
)

type respondent struct {
	replacer   Replacer
	expander   Expander
	applicator Applicator
}

type HttpContext struct {
	W http.ResponseWriter
	R *http.Request
}

var ErrBadExpander = errors.New("respondent: expander is required")

func (rs *respondent) Callback(ctx any, err error) {
	errRep := rs.replacer.Replace(err)
	if errRep == nil {
		return
	}
	manifest := rs.expander.Expand(errRep)
	if manifest == nil {
		return
	}
	rs.applicator.Apply(ctx, manifest)
}

func (rs *respondent) CallbackForHttp(w http.ResponseWriter, r *http.Request, err error) {
	ctxHttp := &HttpContext{
		W: w,
		R: r,
	}
	rs.Callback(ctxHttp, err)
}

func newRespondent(expander Expander, replacer Replacer, applicator Applicator) *respondent {
	if expander == nil {
		panic(ErrBadExpander)
	}
	if replacer == nil {
		replacer = NewSimpleReplacer()
	}
	if applicator == nil {
		applicator = NewSimpleApplicator()
	}

	resp := &respondent{
		expander:   expander,
		replacer:   replacer,
		applicator: applicator,
	}
	return resp
}

func NewMiddleware(expander Expander, replacer Replacer, applicator Applicator) httph.Middleware {
	resp := newRespondent(expander, replacer, applicator)
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(w, r)
			err := httph.ErrorGet(r)
			boo := httph.ErrorTryAcquireHandling(r)
			if err != nil && boo {
				resp.CallbackForHttp(w, r, err)
			}
		})
	}
}
