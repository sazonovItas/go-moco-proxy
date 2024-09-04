package cmd

import (
	"github.com/sazonovItas/go-moco-proxy/internal/app"
	"github.com/sazonovItas/go-moco-proxy/pkg/logger"
	"github.com/spf13/cobra"
)

const (
	listenerFlag = "listener"
	targetFlag   = "target"
	mirrorFlag   = "mirror"
	metricsFlag  = "metrics"
)

type serveCmd struct {
	cmd  *cobra.Command
	opts serveOpts
}

type serveOpts struct {
	listener string
	targets  []string
	mirror   string
	metrics  string
}

func newServeCmd() *serveCmd {
	root := &serveCmd{}
	cmd := &cobra.Command{
		Use:               "serve",
		Short:             "Run tcp proxy with given listener, targets, mirror and metric addresses.",
		Long:              "Run tcp proxy with given listener, targets, mirror and metric addresses.",
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(_ *cobra.Command, _ []string) error {
			options := root.opts

			cfg, err := generateConfig(
				options.listener,
				options.targets,
				options.mirror,
				options.metrics,
			)
			if err != nil {
				return err
			}

			application, err := app.NewApp(logger.NewLogger(logger.CreateLogger()), cfg)
			if err != nil {
				return err
			}
			defer application.Shutdown()

			return application.Run()
		},
	}

	cmd.PersistentFlags().
		StringVarP(&root.opts.listener, listenerFlag, "l", "", "Specify proxy listen address")
	_ = cmd.MarkFlagRequired(listenerFlag)
	cmd.PersistentFlags().
		StringSliceVarP(&root.opts.targets, "target", "t", nil, "Specify proxy target addresses")
	_ = cmd.MarkFlagRequired(targetFlag)
	cmd.MarkFlagsRequiredTogether(listenerFlag, targetFlag)
	cmd.MarkFlagsOneRequired(listenerFlag, targetFlag)
	cmd.PersistentFlags().
		StringVarP(&root.opts.mirror, mirrorFlag, "m", "", "Specify proxy mirror address")
	cmd.PersistentFlags().
		StringVar(&root.opts.metrics, metricsFlag, "", "Specify proxy metric address")

	root.cmd = cmd
	return root
}
