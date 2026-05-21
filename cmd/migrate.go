package cmd

import (
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/AlexBond702/catalog-service/internal/app/builder"
)

const (
	cmdMigrateUsage = "Применяет миграции базы данных при наличии новых"

	cmdMigrateDescription = ` Устанавливает соединение к Postgres базе данных, проверяет соединение,
и затем применяет те миграции, которые еще не были применены,
в соответствии со схемой данных.
`
)

func Migrate() *cli.Command {
	return &cli.Command{
		Name:            "migrate",
		Aliases:         []string{"m"},
		Usage:           cmdMigrateUsage,
		Description:     strings.TrimSpace(cmdMigrateDescription),
		Action:          cmdMigrate,
		HideHelpCommand: true,
	}
}

func cmdMigrate(cCtx *cli.Context) error {
	app := builder.NewBuilder(cCtx)
	app.BuildConfig()
	app.BuildRepoConnPostgres()
	app.BuildRepoConnMigrate()
	app.Run()
	return nil
}
