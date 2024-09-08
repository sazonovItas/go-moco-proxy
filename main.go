package main

import (
	"fmt"
	"os"

	goversion "github.com/caarlos0/go-version"
	cmd "github.com/sazonovItas/go-moco-proxy/cmd/moco-proxy"
	"github.com/sazonovItas/go-moco-proxy/pkg/logger"
)

//nolint:gochecknoglobals
var (
	version string
	commit  string
	date    string
)

func main() {
	err := logger.ConfigureLogger(
		logger.WithLevel(logger.ParseLevel("info")),
		logger.WithOutputPaths([]string{"stdout"}),
		logger.WithErrorOutputPaths([]string{"stderr"}),
		logger.WithPrettyConsoleEncoding(),
	)
	if err != nil {
		panic(fmt.Errorf("failed to init configure logger: %w", err))
	}
	//nolint:errcheck
	defer logger.Sync()

	cmd.Execute(
		buildVersion(version, commit, date),
		os.Exit,
		os.Args[1:],
	)
}

func buildVersion(version, commit, date string) goversion.Info {
	return goversion.GetVersionInfo(
		goversion.WithAppDetails(
			"moco-proxy",
			"Run TCP/TLS load balancer proxy with mirroring support.",
			"",
		),
		func(i *goversion.Info) {
			if commit != "" {
				i.GitCommit = commit
			}
			if date != "" {
				i.BuildDate = date
			}
			if version != "" {
				i.GitVersion = version
			}
		},
	)
}
