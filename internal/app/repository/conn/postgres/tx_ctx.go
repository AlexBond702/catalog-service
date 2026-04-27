package rcpostgres

import (
	"context"
	"github.com/uptrace/bun"
)

type _ctxKeyTx struct{}

func getTxFromContext(ctx context.Context) bun.Tx {
	value, _ := ctx.Value(_ctxKeyTx{}).(bun.Tx)
	return value

}
func setTxFromContext(ctx context.Context, tx bun.Tx) context.Context {
	ctxValueTx := context.WithValue(ctx, _ctxKeyTx{}, tx)
	return ctxValueTx
}
