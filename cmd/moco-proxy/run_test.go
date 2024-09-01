package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []string
		opts    runOpts
		wantErr bool
	}{
		{
			name:    "no args",
			args:    []string{},
			opts:    runOpts{},
			wantErr: false,
		},
		{
			name: "specify config path",
			args: []string{"-c", "config.yaml"},
			opts: runOpts{
				configPath: "config.yaml",
			},
			wantErr: false,
		},
		{
			name: "specify config path several times",
			args: []string{"-c", "config.yaml", "-c", ".conf"},
			opts: runOpts{
				configPath: ".conf",
			},
			wantErr: false,
		},
		{
			name:    "unknown flag",
			args:    []string{"--xxxxxx"},
			opts:    runOpts{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := newRunCmd()
			cmd := root.cmd
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			require.True(t, (err != nil) == tt.wantErr, "%v", err)
			require.Equal(t, tt.opts, root.opts)
		})
	}
}
