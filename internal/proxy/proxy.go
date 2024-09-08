package proxy

import (
	"context"
	"net"
	"sync"

	"github.com/libp2p/go-reuseport"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sazonovItas/go-moco-proxy/pkg/config"
	"github.com/sazonovItas/go-moco-proxy/pkg/connpool"
	"github.com/sazonovItas/go-moco-proxy/pkg/logger"
	"github.com/sazonovItas/go-moco-proxy/pkg/metrics"
	"go.uber.org/zap"
)

var (
	upstreamConnActive = metrics.MustRegisterGauge(
		"server",
		"upstream_conn_active",
		"count of current active upstream connections",
		"address",
	)
	upstreamConnTotal = metrics.MustRegisterCounter(
		"server",
		"upstream_conn_total",
		"total upstream connections",
		"address",
	)
	upstreamConnErr = metrics.MustRegisterCounter(
		"server",
		"upstream_conn_error",
		"total upstream connection errors",
		"address",
	)

	downstremConnActive = metrics.MustRegisterGauge(
		"server",
		"downstream_conn_active",
		"count of current active downstream connections",
		"name",
	)
	downstremConnTotal = metrics.MustRegisterCounter(
		"server",
		"downstream_conn_total",
		"total downstream connections",
		"name",
	)
	downstremConnErr = metrics.MustRegisterCounter(
		"server",
		"downstream_conn_err",
		"total downstream connection errors",
		"name",
	)
)

type Pool interface {
	Acquire(ctx context.Context) (connpool.Conn, error)
}

type Proxy struct {
	name string
	cfg  config.HostConfig

	pool     Pool
	listener net.Listener
	cancel   context.CancelFunc
	sync.Mutex

	wg sync.WaitGroup
}

// NewProxy functions returns new proxy with given name, host config and logger.
func NewProxy(name string, pool Pool, cfg config.HostConfig) *Proxy {
	return &Proxy{
		name: name,
		pool: pool,
		cfg:  cfg,
	}
}

// ListenAndServe method initializates new listener with reuseport,
// and canceled previous listen context if cancelation function is not nil.
func (p *Proxy) ListenAndServe(ctx context.Context) {
	p.Lock()
	if p.cancel != nil {
		p.cancel()
		p.cancel = nil
	}

	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	p.Unlock()

	listener, err := reuseport.Listen("tcp", p.cfg.Address)
	if err != nil {
		logger.FromContext(ctx).Error(
			"failed to setup listener",
			zap.String("address", p.cfg.Address),
			zap.Error(err),
		)
		return
	}

	p.Lock()
	p.listener = listener
	p.Unlock()

	p.serve(ctx)
}

// serve method accepts new connections from listener.
func (p *Proxy) serve(ctx context.Context) {
	for {
		conn, err := p.listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				downstremConnErr.With(prometheus.Labels{"name": p.name}).Inc()
				logger.FromContext(ctx).Error("failed to serve connection", zap.Error(err))
				continue
			}
		}

		downstremConnTotal.With(prometheus.Labels{"name": p.name}).Inc()
		downstremConnActive.With(prometheus.Labels{"name": p.name}).Inc()

		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			p.handleConn(conn)
			downstremConnActive.With(prometheus.Labels{"name": p.name}).Dec()
		}()
	}
}

// handleConn method handles connection accepted by listener.
func (p *Proxy) handleConn(_ net.Conn) {
	upstreamConnTotal.With(prometheus.Labels{"address": ""}).Inc()
	upstreamConnActive.With(prometheus.Labels{"address": ""}).Inc()
	defer upstreamConnActive.With(prometheus.Labels{"address": ""}).Dec()
}

// Shutdown method cancels context and closes proxy listener.
func (p *Proxy) Shutdown() {
	p.Lock()
	if p.cancel != nil {
		p.cancel()
		p.cancel = nil
	}

	if p.listener != nil {
		p.listener.Close()
	}
	p.wg.Wait()
	p.Unlock()
}
