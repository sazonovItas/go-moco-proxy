package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"github.com/sazonovItas/go-moco-proxy/pkg/config"
)

// loadConfig loads config from configPath or other default paths
// and returns path to config that was using to load config.
func loadConfig(configPath string) (cfg *config.Config, used string, err error) {
	if configPath != "" {
		cfg, err = config.Load(configPath)
		if err != nil {
			return nil, "", err
		}

		return cfg, configPath, nil
	}

	for _, path := range [8]string{
		"moco-proxy.yaml",
		"moco-proxy.yml",
		".moco-proxy.yaml",
		".moco-proxy.yml",
		"$HOME/.moco-proxy/config.yaml",
		"$HOME/.moco-proxy/config.yml",
		"/etc/moco-proxy/config.yaml",
		"/etc/moco-proxy/config.yml",
	} {
		cfg, err = config.Load(path)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}

			return nil, "", fmt.Errorf("failed to load config from %s: %w", path, err)
		}

		return cfg, path, err
	}

	return nil, "", fmt.Errorf("failed to load config from default paths")
}

// generateConfig generates config by given listener, targets, mirror and metrics.
func generateConfig(
	listener string,
	targets []string,
	mirror string,
	metrics string,
) (cfg *config.Config, err error) {
	l, err := parseHostConfigFromAddr(listener)
	if err != nil {
		return nil, err
	}

	server := config.ServerConfig{
		Name:     "default",
		Listener: l,
	}

	for _, target := range targets {
		t, err := parseHostConfigFromAddr(target)
		if err != nil {
			return nil, err
		}

		server.Targets = append(server.Targets, t)
	}

	if len(mirror) != 0 {
		m, err := parseHostConfigFromAddr(mirror)
		if err != nil {
			return nil, err
		}

		server.Mirror = m
	}

	cfg = &config.Config{
		Servers: []config.ServerConfig{server},
	}
	if len(metrics) != 0 {
		m, err := parseHostConfigFromAddr(metrics)
		if err != nil {
			return nil, err
		}

		cfg.Metrics = config.MetricConfig{Address: m.Address}
	}

	return cfg, nil
}

func parseHostConfigFromAddr(addr string) (config.HostConfig, error) {
	_, err := parseAddr(addr)
	if err != nil {
		return config.HostConfig{}, err
	}

	return config.HostConfig{Address: addr}, nil
}

func parseAddr(addr string) (hostPort []string, err error) {
	hostPort = strings.Split(addr, ":")
	if len(hostPort) != 2 {
		return nil, fmt.Errorf("invalid address format %s", addr)
	}

	return hostPort, nil
}
