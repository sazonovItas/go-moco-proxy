package cmd

import (
	"reflect"
	"testing"

	"github.com/sazonovItas/go-moco-proxy/pkg/config"
	"github.com/stretchr/testify/require"
)

func Test_loadConfig(t *testing.T) {
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
		// TODO: Add test cases.
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

			require.Equal(t, tt.want, got, "generateConfig() got = %v, want %v", got, tt.want)
		})
	}
}
