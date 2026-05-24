package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/AlexBond702/catalog-service/cmd"
)

func main() {
	app := &cli.App{
		Name:  "Catalog-Service",
		Usage: "Catalog management service",
		Commands: []*cli.Command{
			cmd.Migrate(),
			cmd.WebServer(),
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "no-json"},
		},
		Version: "1.00.00",
	}
	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
