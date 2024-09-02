package cmd

import (
	"reflect"
	"testing"

	"github.com/sazonovItas/go-moco-proxy/internal/config"
	"github.com/stretchr/testify/require"
)

func Test_generateConfig(t *testing.T) {
	t.Parallel()

	type args struct {
		listener string
		targets  []string
		mirror   string
		metrics  string
	}
	tests := []struct {
		name    string
		args    args
		want    *config.Config
		wantErr bool
	}{
		{
			name: "valid config with listener and targets",
			args: args{
				listener: "127.0.0.1:8080",
				targets: []string{
					"127.0.0.1:10000",
					"127.0.0.1:10001",
					"127.0.0.1:10002",
				},
			},
			want: &config.Config{
				Servers: []config.ServerConfig{
					{
						Name: "default",
						Listener: config.HostConfig{
							Host: "127.0.0.1",
							Port: "8080",
						},
						Targets: []config.HostConfig{
							{
								Host: "127.0.0.1",
								Port: "10000",
							},
							{
								Host: "127.0.0.1",
								Port: "10001",
							},
							{
								Host: "127.0.0.1",
								Port: "10002",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid config with listener, targets, mirror and metrics",
			args: args{
				listener: "127.0.0.1:8080",
				targets: []string{
					"127.0.0.1:10000",
				},
				mirror:  "127.0.0.1:3030",
				metrics: "127.0.0.1:4040",
			},
			want: &config.Config{
				Servers: []config.ServerConfig{
					{
						Name: "default",
						Listener: config.HostConfig{
							Host: "127.0.0.1",
							Port: "8080",
						},
						Targets: []config.HostConfig{
							{
								Host: "127.0.0.1",
								Port: "10000",
							},
						},
						Mirror: config.HostConfig{
							Host: "127.0.0.1",
							Port: "3030",
						},
					},
				},
				Metrics: config.HostConfig{
					Host: "127.0.0.1",
					Port: "4040",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid listener address",
			args: args{
				listener: "8080",
			},
			wantErr: true,
		},
		{
			name: "invalid target address",
			args: args{
				listener: "127.0.0.1:8080",
				targets:  []string{""},
			},
			wantErr: true,
		},
		{
			name: "invalid mirror address",
			args: args{
				listener: "127.0.0.1:8080",
				targets:  []string{"127.0.0.1:8080"},
				mirror:   "8080",
			},
			wantErr: true,
		},
		{
			name: "invalid metrics address",
			args: args{
				listener: "127.0.0.1:8080",
				targets:  []string{"127.0.0.1:8080"},
				mirror:   "127.0.0.1:8080",
				metrics:  "8080",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateConfig(
				tt.args.listener,
				tt.args.targets,
				tt.args.mirror,
				tt.args.metrics,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			require.Equal(t, tt.want, got, "generateConfig() got = %v, want %v", got, tt.want)
		})
	}
}

func Test_loadConfig(t *testing.T) {
	t.Parallel()

	type args struct {
		configPath string
	}
	tests := []struct {
		name     string
		args     args
		wantCfg  *config.Config
		wantUsed string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCfg, gotUsed, err := loadConfig(tt.args.configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCfg, tt.wantCfg) {
				t.Errorf("loadConfig() gotCfg = %v, want %v", gotCfg, tt.wantCfg)
			}
			if gotUsed != tt.wantUsed {
				t.Errorf("loadConfig() gotUsed = %v, want %v", gotUsed, tt.wantUsed)
			}
		})
	}
}
