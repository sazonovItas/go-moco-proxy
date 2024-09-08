package cmd

import (
	"fmt"

	"github.com/sazonovItas/go-moco-proxy/internal/app"
	"github.com/sazonovItas/go-moco-proxy/pkg/logger"
	"github.com/spf13/cobra"
)

type runCmd struct {
	cmd  *cobra.Command
	opts runOpts
}

type runOpts struct {
	configPath string
}

func newRunCmd() *runCmd {
	root := &runCmd{}
	cmd := &cobra.Command{
		Use:               "run",
		Short:             "Run proxy with given config.",
		Long:              "Run proxy with given config.",
		Aliases:           []string{"r"},
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg, path, err := loadConfig(root.opts.configPath)
			if err != nil {
				return err
			}
			logger.Info(fmt.Sprintf("Using %s config file", path))

			application, err := app.NewApp(cfg)
			if err != nil {
				return err
			}
			defer application.Shutdown()

			return application.Run()
		},
	}

	cmd.PersistentFlags().
		StringVarP(&root.opts.configPath, "config", "c", "", "Specify path to config file")
	_ = cmd.MarkFlagFilename("config", "yaml", "yml")
	_ = cmd.MarkFlagRequired("config")

	root.cmd = cmd
	return root
}
