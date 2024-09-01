package cmd

import (
	goversion "github.com/caarlos0/go-version"
	"github.com/sazonovItas/go-moco-proxy/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type rootCmd struct {
	cmd  *cobra.Command
	opts rootOpts
	exit func(int)
}

type rootOpts struct {
	debug bool
}

func Execute(version goversion.Info, exit func(int), args []string) {
	newRootCmd(version, exit).Execute(args)
}

func (cmd *rootCmd) Execute(args []string) {
	cmd.cmd.SetArgs(args)

	if err := cmd.cmd.Execute(); err != nil {
		logger.Error("command failed", zap.Error(err))
		cmd.exit(1)
	}
}

func newRootCmd(version goversion.Info, exit func(int)) *rootCmd {
	root := &rootCmd{
		exit: exit,
	}

	cmd := &cobra.Command{
		Use:   "moco-proxy [command]",
		Short: "Run TCP/TLS load balancer proxy with mirroring support.",
		Long: `Run TCP/TLS load balancer proxy with mirroring support. 
You can specify prometheus metrics server for monitoring proxy.`,
		Version:           version.String(),
		SilenceUsage:      true,
		SilenceErrors:     false,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) (err error) {
			if root.opts.debug {
				if err = logger.ConfigureLogger(
					logger.WithDevelopmentLogs(true),
					logger.WithLevel(logger.ParseLevel("debug")),
				); err != nil {
					return
				}
			}

			return
		},
	}
	cmd.SetVersionTemplate("{{ .Version }}")
	cmd.PersistentFlags().
		BoolVarP(&root.opts.debug, "debug-logs", "D", false, "Enable debug logs")

	cmd.AddCommand(
		newRunCmd().cmd,
		newServeCmd().cmd,
	)

	root.cmd = cmd
	return root
}
