package main

import (
	"context"
	"os/signal"
	"syscall"

	vinfo "github.com/sazonovItas/go-moco-proxy/pkg/version"
)

var (
	version  string
	commit   string
	branch   string
	date     string
	platform string
)

func main() {
	_ = vinfo.Info{
		GitVersion: version,
		GitCommit:  commit,
		GitBranch:  branch,
		BuildDate:  date,
		Platform:   platform,
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	defer cancel()

	<-ctx.Done()
}
