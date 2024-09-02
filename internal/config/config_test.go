package config

import (
	"testing"
)

func TestTLSConfig_IsSimple(t *testing.T) {
	tests := []struct {
		name string
		mode string
		want bool
	}{
		{
			name: "simple tls mode",
			mode: SimpleTLSMode,
			want: true,
		},
		{
			name: "insecure tls mode",
			mode: InsecureTLSMode,
			want: false,
		},
		{
			name: "mutual tls mode",
			mode: MutualTLSMode,
			want: false,
		},
		{
			name: "unknown tls mode",
			mode: "unknown",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := TLSConfig{
				Mode: tt.mode,
			}
			if got := tr.IsSimple(); got != tt.want {
				t.Errorf("TLSConfig.IsSimple() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTLSConfig_IsInsecure(t *testing.T) {
	tests := []struct {
		name string
		mode string
		want bool
	}{
		{
			name: "simple tls mode",
			mode: SimpleTLSMode,
			want: false,
		},
		{
			name: "insecure tls mode",
			mode: InsecureTLSMode,
			want: true,
		},
		{
			name: "mutual tls mode",
			mode: MutualTLSMode,
			want: false,
		},
		{
			name: "unknown tls mode",
			mode: "unknown",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := TLSConfig{
				Mode: tt.mode,
			}
			if got := tr.IsInsecure(); got != tt.want {
				t.Errorf("TLSConfig.IsInsecure() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTLSConfig_IsMutual(t *testing.T) {
	tests := []struct {
		name string
		mode string
		want bool
	}{
		{
			name: "simple tls mode",
			mode: SimpleTLSMode,
			want: false,
		},
		{
			name: "insecure tls mode",
			mode: InsecureTLSMode,
			want: false,
		},
		{
			name: "mutual tls mode",
			mode: MutualTLSMode,
			want: true,
		},
		{
			name: "unknown tls mode",
			mode: "unknown",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := TLSConfig{
				Mode: tt.mode,
			}
			if got := tr.IsMutual(); got != tt.want {
				t.Errorf("TLSConfig.IsMutual() = %v, want %v", got, tt.want)
			}
		})
	}
}
