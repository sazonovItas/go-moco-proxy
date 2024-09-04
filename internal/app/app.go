package app

import (
	"context"

	"github.com/sazonovItas/go-moco-proxy/pkg/config"
	"github.com/sazonovItas/go-moco-proxy/pkg/logger"
	"github.com/urfave/cli/v2"
)

type App interface {
	Context() context.Context
	Logger() logger.Logger
	Run()
	Shutdown()
}

type app struct {
	cliCommand *cli.Command
	ctx        context.Context
	cancel     context.CancelFunc

	log logger.Logger
	cfg *config.Config
}

func NewApp(log logger.Logger, cfg *config.Config) (*app, error) {
	return &app{
		log: log,
		cfg: cfg,
	}, nil
}

func (a *app) Run() error {
	return nil
}

func (a *app) Context() context.Context {
	return a.ctx
}

func (a *app) Logger() logger.Logger {
	return a.log
}

func (a *app) Shutdown() {
	a.log.Info("Shutdown proxy app")
	if a.cancel != nil {
		a.cancel()
	}
}
