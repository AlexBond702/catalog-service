package cmd

import (
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/AlexBond702/catalog-service/internal/app/builder"
)

const (
	cmdWebServerUsage = "Starts the web (REST) server"

	cmdWebServerDescription = `
Initializes and starts web-server, that listens specified port
for incoming REST requests.
`
)

func WebServer() *cli.Command {
	return &cli.Command{
		Name:            "web-server",
		Aliases:         []string{"web", "http"},
		Usage:           cmdWebServerUsage,
		Description:     strings.TrimSpace(cmdWebServerDescription),
		Action:          cmdWebServer,
		HideHelpCommand: true,
	}
}

func cmdWebServer(cCtx *cli.Context) error {
	app := builder.NewBuilder(cCtx)
	app.BuildConfig()
	app.BuildRepoConnPostgres()

	app.BuildRepoCategory()
	app.BuildRepoProduct()

	app.BuildServiceCategory()
	app.BuildServiceProduct()

	app.BuildHandlerHttpCategory()
	app.BuildHandlerHttpProduct()

	app.BuildHandlerGrpcCatalog()

	app.BuildProcGrpcGateway() // GATEWAY :8081
	app.BuildProcGrpc()        // GRPC API :50051
	app.BuildProcHttp()        // REST API :8080

	app.Run()
	return nil
}
