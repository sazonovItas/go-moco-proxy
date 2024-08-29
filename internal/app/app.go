package app

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sazonovItas/go-moco-proxy/pkg/logger"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

type Spec struct {
	Name                    string
	Commands                cli.Commands
	IgnoreMissingConfigFile bool
	Defaults                map[string]any
	ConfigDecodingOptions   []viper.DecoderConfigOption
	LateInitTasks           []func(ctx *cli.Context) error
}

type App interface {
	Context() context.Context
	Logger() logger.Logger
	Cli() *cli.App
	MustRun()
	Shutdown()
}

type app struct {
	cliApp *cli.App
	ctx    context.Context
	cancel context.CancelFunc

	log logger.Logger
}

func (a *app) MustRun() {
	if err := a.cliApp.RunContext(a.ctx, os.Args); err != nil {
		if a.log != nil {
			a.log.Error("failed to run proxy", zap.Error(err))
			os.Exit(1)
		} else {
			panic(err)
		}
	}
}

func (a *app) Context() context.Context {
	return a.ctx
}

func (a *app) Logger() logger.Logger {
	return a.log
}

func (a *app) Cli() *cli.App {
	return a.cliApp
}

func (a *app) Shutdown() {
	a.log.Info("Shutdown proxy")
	a.cancel()
}

func NewApp(spec Spec) *app {
	if spec.Defaults == nil {
		spec.Defaults = make(map[string]any)
	}

	lateInitTasks := spec.LateInitTasks
	runCmd := &cli.Command{
		Name:  "run",
		Usage: "Runs proxy with given config",
		Flags: []cli.Flag{},
		Before: func(ctx *cli.Context) (err error) {
			for _, initTask := range lateInitTasks {
				if err = initTask(ctx); err != nil {
					return
				}
			}

			return
		},
	}

	commands := append([]*cli.Command{runCmd}, spec.Commands...)
	return &app{
		cliApp: &cli.App{
			Name:           spec.Name,
			Usage:          "Runs TCP/TLS proxy",
			DefaultCommand: "run",
			Commands:       commands,
		},
	}
}

func (s Spec) readConfig() error {
	viperCfg := viper.NewWithOptions()
	viperCfg.SetConfigName("config")
	viperCfg.SetConfigType("yaml")
	viperCfg.AddConfigPath(fmt.Sprintf("/etc/%s/", strings.ToLower(s.Name)))
	viperCfg.AddConfigPath(fmt.Sprintf("$HOME/.%s/", strings.ToLower(s.Name)))
	viperCfg.AddConfigPath(".")
	viperCfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viperCfg.SetEnvPrefix("MOCO_PROXY")
	viperCfg.AutomaticEnv()

	for k, v := range s.Defaults {
		viperCfg.SetDefault(k, v)
	}

	return nil
}
