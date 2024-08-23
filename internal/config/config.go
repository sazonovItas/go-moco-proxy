package config

import "time"

type Config struct {
	ServerConfigs []ServerConfig `yaml:"servers" koanf:"servers"`
	MetricConfig  *HostConfig    `yaml:"metrics" koanf:"metrics"`
}

type ServerConfig struct {
	Name     string         `yaml:"name"     koanf:"name"`
	Listener *HostConfig    `yaml:"listener" koanf:"listener"`
	Targets  []TargetConfig `yaml:"targets"  koanf:"targets"`
	Mirror   *HostConfig    `yaml:"mirror"   koanf:"mirror"`
	ConnectionConfig
}

type TargetConfig struct {
	Host      string `yaml:"host"       koanf:"host"`
	Port      string `yaml:"port"       koanf:"port"`
	TLSConfig `       yaml:"tls_config" koanf:"tls"`
}

type HostConfig struct {
	Host             string `yaml:"host"              koanf:"host"`
	Port             string `yaml:"port"              koanf:"port"`
	ConnectionConfig `       yaml:"connection_config" koanf:"connection"`
	TLSConfig        `       yaml:"tls_config"        koanf:"tls"`
}

type ConnectionConfig struct {
	ConnectTimeout time.Duration `yaml:"connect_timeout" koanf:"connect_timeout"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"    koanf:"idle_timeout"`
	Timeout        time.Duration `yaml:"timeout"         koanf:"timeout"`
}

type TLSConfig struct {
	CaCert string `yaml:"ca_cert" koanf:"ca_cert"`
	Cert   string `yaml:"cert"    koanf:"cert"`
	Key    string `yaml:"key"     koanf:"key"`
}
