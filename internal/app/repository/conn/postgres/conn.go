package rcpostgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"

	"github.com/AlexBond702/catalog-service/internal/app/config/section"
	"github.com/AlexBond702/catalog-service/migration"

	_ "github.com/AlexBond702/catalog-service/internal/app/util"
)

type (
	Client struct {
		_bunDB
		rawBunDB *bun.DB
		cfg      section.RepositoryPostgres
	}
	_bunDB = bun.IDB
)

func (c *Client) GetRawBunDb() *bun.DB {
	return c.rawBunDB
}

func NewConn(ctx context.Context, cfg section.RepositoryPostgres) (*Client, error) {
	var u url.URL
	u.Scheme = "postgres"
	u.Host = cfg.Address
	u.Path = cfg.Name
	u.User = url.UserPassword(cfg.Username, cfg.Password)
	args := make(url.Values)
	args.Set("sslmode", "disable")
	u.RawQuery = args.Encode()
	dsn := u.String()
	log.Printf("ReadTimeout: %v", cfg.ReadTimeout)
	log.Printf("WriteTimeout: %v", cfg.WriteTimeout)

	sqlConnect := pgdriver.NewConnector(
		pgdriver.WithDSN(dsn),
		pgdriver.WithReadTimeout(cfg.ReadTimeout.Duration),
		pgdriver.WithWriteTimeout(cfg.WriteTimeout.Duration),
	)
	sqlDB := sql.OpenDB(sqlConnect)
	sqlDB.SetMaxOpenConns(10)
	bunDB := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil,
			fmt.Errorf("failed connection ping: %w", err)
	}
	c := Client{
		rawBunDB: bunDB,
		cfg:      cfg,
	}
	c._bunDB = bunDB
	return &c, nil
}

func (c *Client) Migrate(ctx context.Context) (oldVer, newVer int64, err error) {
	migrations := migrate.NewMigrations()
	if err := migrations.Discover(migration.Postgres); err != nil {
		return 0, 0, fmt.Errorf("failed to discover migrations: %w", err)
	}

	migrator := migrate.NewMigrator(
		c.rawBunDB,
		migrations,
		migrate.WithTableName(c.cfg.MigrationTable),
		migrate.WithLocksTableName(c.cfg.MigrationTable+"_lock"),
		migrate.WithMarkAppliedOnSuccess(true),
	)
	if err := migrator.Init(ctx); err != nil {
		return 0, 0, fmt.Errorf("migration init failed: %w", err)
	}

	applied, err := migrator.AppliedMigrations(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("migration apply failed: %w", err)
	}

	for i, m := range applied {
		verseOld, parseErr := strconv.ParseInt(m.Name, 10, 64)
		if parseErr != nil {
			return 0, 0, fmt.Errorf("parsing verseOld %s failed %w", m.Name, parseErr)
		}
		if i == 0 {
			oldVer = verseOld
		}
	}

	newVer = oldVer

	unapply, err := migrator.Migrate(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("migrate failed: %w", err)
	}

	for _, m := range unapply.Migrations {
		verseNew, parseErr := strconv.ParseInt(m.Name, 10, 64)
		if parseErr != nil {
			return 0, 0, fmt.Errorf("parsing verseNew %s failed %w", m.Name, parseErr)
		}
		newVer = verseNew
	}
	return oldVer, newVer, nil
}
