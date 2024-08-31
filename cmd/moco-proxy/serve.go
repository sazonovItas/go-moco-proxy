package cmd

import (
	"github.com/spf13/cobra"
)

type serveCmd struct {
	cmd  *cobra.Command
	opts serveOpts
}

type serveOpts struct {
	listener string
	targets  []string
	mirror   string
	metric   string
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
		Run: func(_ *cobra.Command, _ []string) {
			// TODO: add serve logic
		},
	}

	cmd.PersistentFlags().
		StringVarP(&root.opts.listener, "listener", "l", "", "Specify proxy listen address")
	_ = cmd.MarkFlagRequired("listener")
	cmd.PersistentFlags().
		StringSliceVarP(&root.opts.targets, "target", "t", nil, "Specify proxy target addresses")
	_ = cmd.MarkFlagRequired("target")
	cmd.MarkFlagsRequiredTogether("listener", "target")
	cmd.MarkFlagsOneRequired("listener", "target")
	cmd.PersistentFlags().
		StringVarP(&root.opts.mirror, "mirror", "m", "", "Specify proxy mirror address")
	cmd.PersistentFlags().
		StringVar(&root.opts.metric, "metric", "", "Specify proxy metric address")

	root.cmd = cmd
	return root
}
