package config

const (
	SimpleTLSMode   = "simple"
	InsecureTLSMode = "insecure"
	MutualTLSMode   = "mutual"
)

type Config struct {
	Servers []ServerConfig `yaml:"servers" koanf:"servers"`
	Metrics MetricConfig   `yaml:"metrics" koanf:"metrics"`
}

type ServerConfig struct {
	Name     string       `yaml:"name"     koanf:"name"`
	Listener HostConfig   `yaml:"listener" koanf:"listener"`
	Targets  []HostConfig `yaml:"targets"  koanf:"targets"`
	Mirror   HostConfig   `yaml:"mirror"   koanf:"mirror"`
}

type HostConfig struct {
	Address   string `yaml:"address" koanf:"address"`
	TLSConfig `       yaml:"tls"     koanf:"tls"`
}

type TLSConfig struct {
	CaCert string `yaml:"ca_cert" koanf:"ca_cert"`
	Cert   string `yaml:"cert"    koanf:"cert"`
	Key    string `yaml:"key"     koanf:"key"`
	Mode   string `yaml:"mode"    koanf:"mode"`
}

type MetricConfig struct {
	Address   string `yaml:"address"  koanf:"address"`
	Endpoint  string `yaml:"endpoint" koanf:"endpoint"`
	TLSConfig `       yaml:"tls"      koanf:"tls"`
}

func (t TLSConfig) IsSimple() bool {
	return t.Mode == SimpleTLSMode
}

func (t TLSConfig) IsInsecure() bool {
	return t.Mode == InsecureTLSMode
}

func (t TLSConfig) IsMutual() bool {
	return t.Mode == MutualTLSMode
}
