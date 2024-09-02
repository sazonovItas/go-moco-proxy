package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %s: %w", path, err)
	}
	defer f.Close()

	return LoadReader(f)
}

func LoadReader(fd io.Reader) (cfg *Config, err error) {
	content, err := io.ReadAll(fd)
	if err != nil {
		return nil, fmt.Errorf("failed to read config from file: %w", err)
	}

	if err = yaml.Unmarshal(content, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return
}
