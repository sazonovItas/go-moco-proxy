package cmd

import (
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
		Run: func(_ *cobra.Command, _ []string) {
			// TODO: add run logic
		},
	}

	cmd.PersistentFlags().
		StringVarP(&root.opts.configPath, "config", "c", "", "Specify path to config file")
	_ = cmd.MarkFlagFilename("config", "yaml", "yml")
	_ = cmd.MarkFlagRequired("config")
	cmd.MarkFlagsOneRequired("config")

	root.cmd = cmd
	return root
}
