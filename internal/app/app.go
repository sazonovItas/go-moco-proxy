package app

import (
	"context"

	"github.com/sazonovItas/go-moco-proxy/pkg/config"
	"github.com/sazonovItas/go-moco-proxy/pkg/logger"
)

type app struct {
	ctx    context.Context
	cancel context.CancelFunc

	cfg *config.Config
}

func NewApp(cfg *config.Config) (*app, error) {
	return &app{
		cfg: cfg,
	}, nil
}

func (a *app) Run() error {
	return nil
}

func (a *app) Context() context.Context {
	return a.ctx
}

func (a *app) Shutdown() {
	logger.Info("Shutdown proxy app")
	if a.cancel != nil {
		a.cancel()
	}
}
