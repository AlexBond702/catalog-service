package rcpostgres

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

type implBunIdbTxInjector struct {
	fallback bun.IDB
	sqlDB    *sql.DB
}

func newBunIdbTxInjector(orig bun.IDB, sqlDB *sql.DB) bun.IDB {
	return &implBunIdbTxInjector{
		fallback: orig,
		sqlDB:    sqlDB,
	}
}

type sqlConn interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func (x *implBunIdbTxInjector) getIdb(ctx context.Context) bun.IDB {
	tx := getTxFromContext(ctx)
	if tx.Tx != nil {
		return tx
	}
	return x.fallback
}

func (x *implBunIdbTxInjector) raw(ctx context.Context) sqlConn {
	tx := getTxFromContext(ctx)
	if tx.Tx != nil {
		return tx.Tx
	}
	return x.sqlDB
}

func (x *implBunIdbTxInjector) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	getTx := x.raw(ctx)
	return getTx.QueryContext(ctx, query, args...)
}

func (x *implBunIdbTxInjector) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	getTx := x.raw(ctx)
	return getTx.ExecContext(ctx, query, args...)
}

func (x *implBunIdbTxInjector) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	getTx := x.raw(ctx)
	return getTx.QueryRowContext(ctx, query, args...)
}

func (x *implBunIdbTxInjector) NewSelect() *bun.SelectQuery {
	query := x.fallback.NewSelect()
	q := query.Conn(x)
	return q
}

func (x *implBunIdbTxInjector) NewInsert() *bun.InsertQuery {
	query := x.fallback.NewInsert()
	q := query.Conn(x)
	return q
}

func (x *implBunIdbTxInjector) NewUpdate() *bun.UpdateQuery {
	query := x.fallback.NewUpdate()
	q := query.Conn(x)
	return q
}

func (x *implBunIdbTxInjector) NewDelete() *bun.DeleteQuery {
	query := x.fallback.NewDelete()
	q := query.Conn(x)
	return q
}

func (x *implBunIdbTxInjector) Dialect() schema.Dialect {
	return x.fallback.Dialect()
}

func (x *implBunIdbTxInjector) NewRaw(query string, args ...interface{}) *bun.RawQuery {
	newRaw := x.fallback.NewRaw(query, args)
	conn := newRaw.Conn(x)
	return conn
}

func (x *implBunIdbTxInjector) NewValues(model interface{}) *bun.ValuesQuery {
	newValues := x.fallback.NewValues(model)
	return newValues.Conn(x)
}

func (x *implBunIdbTxInjector) BeginTx(ctx context.Context, opts *sql.TxOptions) (bun.Tx, error) {
	getTx := x.getIdb(ctx)
	get, _ := getTx.BeginTx(ctx, opts)
	return get, nil
}

func (x *implBunIdbTxInjector) RunInTx(ctx context.Context,
	opts *sql.TxOptions,
	f func(ctx context.Context, tx bun.Tx) error,
) error {
	getTx := x.getIdb(ctx)
	return getTx.RunInTx(ctx, opts, f)
}

func (x *implBunIdbTxInjector) NewCreateTable() *bun.CreateTableQuery {
	return x.fallback.NewCreateTable().Conn(x)
}

func (x *implBunIdbTxInjector) NewDropTable() *bun.DropTableQuery {
	return x.fallback.NewDropTable().Conn(x)
}

func (x *implBunIdbTxInjector) NewCreateIndex() *bun.CreateIndexQuery {
	return x.fallback.NewCreateIndex().Conn(x)
}

func (x *implBunIdbTxInjector) NewTruncateTable() *bun.TruncateTableQuery {
	return x.fallback.NewTruncateTable().Conn(x)
}

func (x *implBunIdbTxInjector) NewAddColumn() *bun.AddColumnQuery {
	return x.fallback.NewAddColumn().Conn(x)
}

func (x *implBunIdbTxInjector) NewDropColumn() *bun.DropColumnQuery {
	return x.fallback.NewDropColumn().Conn(x)
}

func (x *implBunIdbTxInjector) NewMerge() *bun.MergeQuery {
	return x.fallback.NewMerge().Conn(x)
}

func (x *implBunIdbTxInjector) NewDropIndex() *bun.DropIndexQuery {
	return x.fallback.NewDropIndex().Conn(x)
}
