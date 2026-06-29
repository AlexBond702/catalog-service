package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/AlexBond702/catalog-service/cmd"
	"github.com/AlexBond702/catalog-service/internal/app/constant"
	msentry "github.com/AlexBond702/catalog-service/internal/app/processor/monitor/sentry"
)

func main() {
	defer msentry.Flush()
	app := &cli.App{
		Name:  constant.AppName,
		Usage: "Catalog management service",
		Commands: []*cli.Command{
			cmd.Migrate(),
			cmd.WebServer(),
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "no-json"},
		},
		Version: constant.Version,
	}
	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
