package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO: change ip addresses in serve command.
func TestServeCmd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []string
		opts    serveOpts
		wantErr bool
	}{
		{
			name:    "no args",
			args:    []string{},
			opts:    serveOpts{},
			wantErr: true,
		},
		{
			name: "specify listener and targets",
			args: []string{"-l", "XXXX:XX", "-t=XXXX:XX", "-t=XXXX:XX"},
			opts: serveOpts{
				listener: "XXXX:XX",
				targets:  []string{"XXXX:XX", "XXXX:XX"},
			},
			wantErr: false,
		},
		{
			name: "specify only targets",
			args: []string{"-t=XXXX:XX", "-t=XXXX:XX"},
			opts: serveOpts{
				targets: []string{"XXXX:XX", "XXXX:XX"},
			},
			wantErr: true,
		},
		{
			name: "specify only listener",
			args: []string{"-l", "XXXX:XX"},
			opts: serveOpts{
				listener: "XXXX:XX",
			},
			wantErr: true,
		},
		{
			name: "specify all flags",
			args: []string{
				"-l",
				"XXXX:XX",
				"-t=XXXX:XX",
				"-t=XXXX:XX",
				"-m=XXXX:XX",
				"--metric",
				"XXXX:XX",
			},
			opts: serveOpts{
				listener: "XXXX:XX",
				targets:  []string{"XXXX:XX", "XXXX:XX"},
				mirror:   "XXXX:XX",
				metric:   "XXXX:XX",
			},
			wantErr: false,
		},
		{
			name: "unknown flag",
			args: []string{
				"--xxxxx",
			},
			opts:    serveOpts{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := newServeCmd()
			cmd := root.cmd
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			require.True(t, (err != nil) == tt.wantErr, "%v", err)
			require.Equal(t, tt.opts, root.opts)
		})
	}
}
