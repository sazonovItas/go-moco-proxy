package cmd

import (
	goversion "github.com/caarlos0/go-version"
	"github.com/sazonovItas/go-moco-proxy/pkg/logger"
	"github.com/spf13/cobra"
)

type rootCmd struct {
	cmd         *cobra.Command
	development bool
	pretty      bool
	exit        func(int)
}

func newRootCmd(version goversion.Info, exit func(int)) *rootCmd {
	root := &rootCmd{
		exit: exit,
	}

	cmd := &cobra.Command{
		Use:               "moco-proxy",
		Short:             "Run TCP/TLS load balancer proxy with mirroring support",
		Long:              "Run TCP/TLS load balancer proxy with mirroring support. You can specify metrics server config to see proxy metrics",
		Version:           version.String(),
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) (err error) {
			if root.development {
				if err = logger.ConfigureLogger(
					logger.WithLevel(logger.ParseLevel("debug")),
					logger.WithDevelopmentLogs(root.development),
				); err != nil {
					return
				}
			}

			if root.pretty {
				if err = logger.ConfigureLogger(); err != nil {
					return
				}
			}

			return
		},
	}
	cmd.SetVersionTemplate("moco-proxy {{.Version}}")
	cmd.PersistentFlags().
		BoolVarP(&root.development, "development", "D", false, "Enable debug logs")
	cmd.PersistentFlags().
		BoolVarP(&root.development, "pretty", "P", false, "Enable pretty logs")

	cmd.AddCommand()

	root.cmd = cmd
	return root
}
