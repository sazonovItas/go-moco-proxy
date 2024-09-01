package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sazonovItas/go-moco-proxy/internal/config"
)

var ErrAddrInvalidFormat = errors.New("invalid address format")

func GenerateConfig(
	listener string,
	targets []string,
	mirror string,
	metrics string,
) (cfg *config.Config, err error) {
	l, err := separateAddress(listener)
	if err != nil {
		return nil, err
	}

	server := config.ServerConfig{
		Name: "default",
		Listener: config.HostConfig{
			Host: l[0],
			Port: l[1],
		},
	}

	cfg = &config.Config{
		Servers: []config.ServerConfig{server},
	}

	for _, target := range targets {
		addr, err := separateAddress(target)
		if err != nil {
			return nil, err
		}

		cfg.Servers[0].Targets = append(cfg.Servers[0].Targets, config.HostConfig{
			Host: addr[0],
			Port: addr[1],
		})
	}

	if mirror != "" {
		addr, err := separateAddress(mirror)
		if err != nil {
			return nil, err
		}

		server.Mirror = &config.HostConfig{
			Host: addr[0],
			Port: addr[1],
		}
	}

	if metrics != "" {
		addr, err := separateAddress(metrics)
		if err != nil {
			return nil, err
		}

		cfg.Metrics = &config.HostConfig{
			Host: addr[0],
			Port: addr[1],
		}
	}

	return
}

func separateAddress(addr string) ([]string, error) {
	hostPort := strings.Split(addr, ":")
	if len(hostPort) != 2 {
		return nil, fmt.Errorf("%w: %s", ErrAddrInvalidFormat, addr)
	}

	return hostPort, nil
}
