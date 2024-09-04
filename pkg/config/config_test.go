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
			name: "Test simple tls mode",
			mode: SimpleTLSMode,
			want: true,
		},
		{
			name: "Test insecure tls mode",
			mode: InsecureTLSMode,
			want: false,
		},
		{
			name: "Test mutual tls mode",
			mode: MutualTLSMode,
			want: false,
		},
		{
			name: "Test unknown tls mode",
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
			name: "Test simple tls mode",
			mode: SimpleTLSMode,
			want: false,
		},
		{
			name: "Test insecure tls mode",
			mode: InsecureTLSMode,
			want: true,
		},
		{
			name: "Test mutual tls mode",
			mode: MutualTLSMode,
			want: false,
		},
		{
			name: "Test unknown tls mode",
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
			name: "Test simple tls mode",
			mode: SimpleTLSMode,
			want: false,
		},
		{
			name: "Test insecure tls mode",
			mode: InsecureTLSMode,
			want: false,
		},
		{
			name: "Test mutual tls mode",
			mode: MutualTLSMode,
			want: true,
		},
		{
			name: "Test unknown tls mode",
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
