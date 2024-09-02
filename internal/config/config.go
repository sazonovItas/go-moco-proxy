package config

const (
	SimpleTLSMode   = "simple"
	InsecureTLSMode = "insecure"
	MutualTLSMode   = "mutual"
)

type Config struct {
	Servers []ServerConfig `yaml:"servers"`
	Metrics HostConfig     `yaml:"metrics"`
}

type ServerConfig struct {
	Name     string       `yaml:"name"`
	Listener HostConfig   `yaml:"listener"`
	Targets  []HostConfig `yaml:"targets"`
	Mirror   HostConfig   `yaml:"mirror"`
}

type HostConfig struct {
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	TLSConfig `yaml:"tls"`
}

type TLSConfig struct {
	CaCert string `yaml:"ca_cert"`
	Cert   string `yaml:"cert"`
	Key    string `yaml:"key"`
	Mode   string `yaml:"mode"`
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
