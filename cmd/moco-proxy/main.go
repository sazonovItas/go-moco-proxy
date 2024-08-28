package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/sazonovItas/go-moco-proxy/internal/metrics"
	"github.com/sazonovItas/go-moco-proxy/pkg/logger"
)

var (
	version string
	commit  string
	date    string
)

const usage = ``

func main() {
	logger.Info("Start metric server")

	srv, err := metrics.New(metrics.MetricConfig{})
	if err != nil {
		panic(err)
	}

	go func() {
		if err := srv.Run(); err != nil {
			logger.Errorf("failed run server: %s", err.Error())
		}
	}()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)
	defer stop()

	<-ctx.Done()

	if err := srv.Shutdown(); err != nil {
		panic(err)
	}
}
