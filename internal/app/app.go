package app

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/sazonovItas/go-moco-proxy/pkg/logger"
	"github.com/urfave/cli/v2"
)

const (
	configFilePathFlag = "config"
	logLevelFlag       = "logs.level"
	logEncodingFlag    = "logs.encoding"
	logDevelopmentFlag = "logs.development"
)

type Spec struct {
	AppName                 string
	CommandName             string
	IgnoreMissingConfigFile bool
	LateInitTasks           []func(ctx *cli.Context) error
}

type App interface {
	Context() context.Context
	Logger() logger.Logger
	CliCommand() *cli.Command
	Run()
	Shutdown()
}

type app struct {
	cliCommand *cli.Command
	ctx        context.Context
	cancel     context.CancelFunc

	log logger.Logger
	cfg Config
}

//nolint:all
func NewApp(spec Spec) *app {
	a := &app{}

	lateInitTasks := append(spec.LateInitTasks, []func(ctx *cli.Context) error{
		func(ctx *cli.Context) error {
			spec.AppName = ctx.App.Name
			cfg, err := spec.readConfig(ctx)
			if err != nil {
				return fmt.Errorf("failed to read config: %w", err)
			}

			a.cfg = cfg
			return nil
		},
		func(ctx *cli.Context) error {
			err := logger.ConfigureLogger(
				logger.WithLevel(logger.ParseLevel(ctx.String(logLevelFlag))),
				logger.WithDevelopmentLogs(ctx.Bool(logDevelopmentFlag)),
				logger.WithEncoding(ctx.String(logEncodingFlag)),
			)
			if err != nil {
				return fmt.Errorf("failed to init logger: %w", err)
			}

			a.log = logger.NewLogger(logger.CreateLogger())

			return nil
		},
	}...)
	runCmd := &cli.Command{
		Name:  spec.CommandName,
		Usage: "Runs proxy with given config",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    configFilePathFlag,
				Usage:   "specify config path `FILE`",
				Aliases: []string{"c"},
			},
			&cli.StringFlag{
				Name:        logLevelFlag,
				Usage:       "specify log level `info`, debug, warn, error, fatal and panic",
				Value:       "info",
				DefaultText: "default `info` level",
			},
			&cli.StringFlag{
				Name:        logEncodingFlag,
				Usage:       "specify log encoding `console` and json",
				Value:       "console",
				DefaultText: "default `console`",
			},
			&cli.StringFlag{
				Name:        logDevelopmentFlag,
				Usage:       "enable development logs",
				DefaultText: "default `false`",
			},
		},
		Action: func(ctx *cli.Context) error {
			return a.Run()
		},
		Before: func(ctx *cli.Context) (err error) {
			for _, initTask := range lateInitTasks {
				if err = initTask(ctx); err != nil {
					return
				}
			}

			return
		},
	}

	a.cliCommand = runCmd

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	a.ctx, a.cancel = ctx, cancel

	return a
}

func (a *app) Run() error {
	<-a.ctx.Done()
	return nil
}

func (a *app) Context() context.Context {
	return a.ctx
}

func (a *app) Logger() logger.Logger {
	return a.log
}

func (a *app) Cli() *cli.Command {
	return a.cliCommand
}

func (a *app) Shutdown() {
	a.log.Info("Shutdown proxy")
	if a.cancel != nil {
		a.cancel()
	}
}
