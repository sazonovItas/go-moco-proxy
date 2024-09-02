package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
			args: []string{"-l", "127.0.0.1:8080", "-t=127.0.0.1:8080", "-t=127.0.0.1:8080"},
			opts: serveOpts{
				listener: "127.0.0.1:8080",
				targets:  []string{"127.0.0.1:8080", "127.0.0.1:8080"},
			},
			wantErr: false,
		},
		{
			name: "specify only targets",
			args: []string{"-t=127.0.0.1:8080", "-t=127.0.0.1:8080"},
			opts: serveOpts{
				targets: []string{"127.0.0.1:8080", "127.0.0.1:8080"},
			},
			wantErr: true,
		},
		{
			name: "specify only listener",
			args: []string{"-l", "127.0.0.1:8080"},
			opts: serveOpts{
				listener: "127.0.0.1:8080",
			},
			wantErr: true,
		},
		{
			name: "specify all flags",
			args: []string{
				"-l",
				"127.0.0.1:8080",
				"-t=127.0.0.1:8080",
				"-t=127.0.0.1:8080",
				"-m=127.0.0.1:8080",
				"--metrics",
				"127.0.0.1:8080",
			},
			opts: serveOpts{
				listener: "127.0.0.1:8080",
				targets:  []string{"127.0.0.1:8080", "127.0.0.1:8080"},
				mirror:   "127.0.0.1:8080",
				metric:   "127.0.0.1:8080",
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
			if err != nil {
				require.True(t, (err != nil) == tt.wantErr, "%v", err)
				return
			}
			require.Equal(t, tt.opts, root.opts)
		})
	}
}
