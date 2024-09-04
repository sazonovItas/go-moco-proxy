package config

import (
	"fmt"
	"io"
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
)

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %s: %w", path, err)
	}
	defer f.Close()

	cfg, err := LoadReader(f)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from file %s: %w", path, err)
	}

	return cfg, nil
}

func LoadReader(fd io.Reader) (*Config, error) {
	k := koanf.New(".")

	content, err := io.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	if err := k.Load(rawbytes.Provider(content), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return &cfg, nil
}
