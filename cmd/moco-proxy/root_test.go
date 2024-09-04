package cmd

import (
	"bytes"
	"testing"

	goversion "github.com/caarlos0/go-version"
	"github.com/stretchr/testify/require"
)

var testVersion = goversion.Info{
	GitVersion: "v1.0.0",
}

func TestRootCmd(t *testing.T) {
	t.Parallel()

	mem := &testMemExit{}
	Execute(testVersion, mem.Exit, []string{"-h"})
	require.Equal(t, 0, mem.code)

	Execute(testVersion, mem.Exit, []string{"--xxxx"})
	require.Equal(t, 1, mem.code)
}

func TestRootCmdVersion(t *testing.T) {
	t.Parallel()

	var b bytes.Buffer
	mem := &testMemExit{}
	cmd := newRootCmd(testVersion, mem.Exit).cmd
	cmd.SetOut(&b)
	cmd.SetArgs([]string{"-v"})
	require.NoError(t, cmd.Execute())
	require.Contains(t, b.String(), testVersion.GitVersion)
	require.Equal(t, 0, mem.code)
}

func Test_RootCmdExecute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []string
		opts    rootOpts
		wantErr bool
	}{
		{
			name:    "Test no args",
			args:    []string{},
			opts:    rootOpts{},
			wantErr: false,
		},
		{
			name:    "Test help flag -h",
			args:    []string{"-h"},
			opts:    rootOpts{},
			wantErr: false,
		},
		{
			name: "Test pretty-logs and debug-logs flags",
			args: []string{"--debug-logs=true"},
			opts: rootOpts{
				debug: true,
			},
			wantErr: false,
		},
		{
			name: "Test pretty-logs and debug-logs flags in short form",
			args: []string{"-D"},
			opts: rootOpts{
				debug: true,
			},
			wantErr: false,
		},
		{
			name:    "Test unknown flag",
			args:    []string{"-P"},
			opts:    rootOpts{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := &testMemExit{}
			root := newRootCmd(goversion.Info{}, mem.Exit)
			cmd := root.cmd
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			require.True(t, (err != nil) == tt.wantErr, "%v", err)
			require.Equal(t, tt.opts, root.opts)
		})
	}
}
