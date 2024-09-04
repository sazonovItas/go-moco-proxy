package config

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

const testDataDir = "testdata/"

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    *Config
		wantErr bool
	}{
		{
			name: "Test valid config",
			path: testDataDir + "config.yml",
			want: &Config{
				Servers: []ServerConfig{
					{
						Name: "server1",
						Listener: HostConfig{
							Address: "127.0.0.1:8080",
							TLSConfig: TLSConfig{
								Cert: "cert.cert",
								Key:  "key.key",
							},
						},
						Targets: []HostConfig{
							{Address: "127.0.0.1:10001"},
							{Address: "127.0.0.1:10002"},
							{
								Address: "127.0.0.1:10003",
								TLSConfig: TLSConfig{
									CaCert: "ca_cert.pem",
									Mode:   "mutual",
								},
							},
						},
					},
					{
						Name:     "server2",
						Listener: HostConfig{Address: "127.0.0.1:9090"},
						Targets: []HostConfig{
							{Address: "127.0.0.1:11001"},
							{Address: "127.0.0.1:11002"},
						},
						Mirror: HostConfig{
							Address: "127.0.0.1:3030",
							TLSConfig: TLSConfig{
								Mode: "simple",
							},
						},
					},
				},
				Metrics: MetricConfig{
					Address:  "127.0.0.1:4040",
					Endpoint: "metrics",
				},
			},
			wantErr: false,
		},
		{
			name:    "Test empty config",
			path:    testDataDir + "empty.yml",
			want:    &Config{},
			wantErr: false,
		},
		{
			name:    "Test config not exists",
			path:    testDataDir + "notexistence.yml",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadReader(t *testing.T) {
	tests := []struct {
		name    string
		fd      io.Reader
		wantCfg *Config
		wantErr bool
	}{
		{
			name: "Test config servers and metrics",
			fd: bytes.NewBufferString(`
servers:
  - name: "server1"
    listener:
      address: "127.0.0.1:8080"	
    targets:
      - address: "127.0.0.1:10001"
      - address: "127.0.0.1:10002"
      - address: "127.0.0.1:10003"
    mirror:
      address: "127.0.0.1:10004"
  - name: "server2"
    listener:
      address: "127.0.0.1:8081"	
    targets:
      - address: "127.0.0.1:11001"
      - address: "127.0.0.1:11002"
      - address: "127.0.0.1:11003"
    mirror:
      address: "127.0.0.1:11004"
metrics:
  address: "127.0.0.1:9090"
  endpoint: "metrics"
`),
			wantCfg: &Config{
				Servers: []ServerConfig{
					{
						Name:     "server1",
						Listener: HostConfig{Address: "127.0.0.1:8080"},
						Targets: []HostConfig{
							{Address: "127.0.0.1:10001"},
							{Address: "127.0.0.1:10002"},
							{Address: "127.0.0.1:10003"},
						},
						Mirror: HostConfig{Address: "127.0.0.1:10004"},
					},
					{
						Name:     "server2",
						Listener: HostConfig{Address: "127.0.0.1:8081"},
						Targets: []HostConfig{
							{Address: "127.0.0.1:11001"},
							{Address: "127.0.0.1:11002"},
							{Address: "127.0.0.1:11003"},
						},
						Mirror: HostConfig{Address: "127.0.0.1:11004"},
					},
				},
				Metrics: MetricConfig{
					Address:  "127.0.0.1:9090",
					Endpoint: "metrics",
				},
			},
		},
		{
			name: "Test tls config in host config",
			fd: bytes.NewBufferString(`
servers:
  - name: "default"
    listener: 
      address: "127.0.0.1:8080"
      tls:
        ca_cert: "ca_cert.pem"
        cert: "cert.cert"
        key: "key.key"
        mode: "insecure"
`),
			wantCfg: &Config{
				Servers: []ServerConfig{
					{
						Name: "default",
						Listener: HostConfig{
							Address: "127.0.0.1:8080",
							TLSConfig: TLSConfig{
								CaCert: "ca_cert.pem",
								Cert:   "cert.cert",
								Key:    "key.key",
								Mode:   "insecure",
							},
						},
					},
				},
			},
		},
		{
			name:    "Test empty buffer",
			fd:      bytes.NewBufferString(""),
			wantCfg: &Config{},
		},
		{
			name: "Test invalid indent config",
			fd: bytes.NewBufferString(`
servers:
  listener:
      address: "127.0.0.1:8080"	
    targets:
      - address: "127.0.0.1:10001"
      - address: "127.0.0.1:10002"
      - address: "127.0.0.1:10003"
`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCfg, err := LoadReader(tt.fd)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotCfg, tt.wantCfg) {
				t.Errorf("LoadReader() = %v, want %v", gotCfg, tt.wantCfg)
			}
		})
	}
}
