package main

import (
	"os"

	_ "embed"

	goversion "github.com/caarlos0/go-version"
	cmd "github.com/sazonovItas/go-moco-proxy/cmd/moco-proxy"
)

//nolint:gochecknoglobals
var (
	version string
	commit  string
	date    string
)

func main() {
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
			"Run TCP/TLS proxy",
			"",
		),
		// goversion.WithASCIIName(art),
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
