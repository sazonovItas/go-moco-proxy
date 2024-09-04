package metrics

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sazonovItas/go-moco-proxy/pkg/config"
	"golang.org/x/net/http2"
)

const (
	// readTimeout is read timeout for metric's http server.
	readTimeout = 5 * time.Second

	// writeTimmeout is write timeout for metric's http server.
	writeTimeout = 5 * time.Second

	// shutdownTimeout is shutdown timeout for metric's http server.
	shutdownTimeout = 15 * time.Second

	// metricEndpoint is default endpoint for metrics.
	metricEndpoint = "/metrics"
)

type MetricServer struct {
	server *http.Server
}

// New function creates new metric server with given metric config.
func New(c config.MetricConfig) (*MetricServer, error) {
	const op = "metrics.New"

	if c.Endpoint == "" {
		c.Endpoint = metricEndpoint
	}

	handler := http.NewServeMux()
	handler.Handle(c.Endpoint, promhttp.Handler())

	server := &http.Server{
		Addr:         c.Address,
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	if err := http2.ConfigureServer(server, nil); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &MetricServer{
		server: server,
	}, nil
}

// Run method starts http server and it blocks until http server stops.
func (m *MetricServer) Run() error {
	const op = "metrics.Run"

	if err := m.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Shutdown method stops http server.
func (m *MetricServer) Shutdown() error {
	const op = "metrics.Shutdown"

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := m.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("%s: failed to shutdown metric's server: %w", op, err)
	}

	return nil
}
