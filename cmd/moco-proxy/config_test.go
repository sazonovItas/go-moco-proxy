package cmd

import (
	"reflect"
	"testing"

	"github.com/sazonovItas/go-moco-proxy/internal/config"
)

func TestGenerateConfig(t *testing.T) {
	type args struct {
		listener string
		targets  []string
		mirror   string
		metrics  string
	}
	tests := []struct {
		name    string
		args    args
		wantCfg *config.Config
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCfg, err := GenerateConfig(tt.args.listener, tt.args.targets, tt.args.mirror, tt.args.metrics)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCfg, tt.wantCfg) {
				t.Errorf("GenerateConfig() = %v, want %v", gotCfg, tt.wantCfg)
			}
		})
	}
}

func Test_separateAddress(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := separateAddress(tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("separateAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("separateAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
