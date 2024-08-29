package main

import (
	"github.com/sazonovItas/go-moco-proxy/internal/app"
	"github.com/urfave/cli/v2"
)

func main() {
	serverApp := app.NewApp(app.Spec{
		Name:     "moco-proxy",
		Commands: []*cli.Command{cliVersion()},
	})

	serverApp.MustRun()
}
