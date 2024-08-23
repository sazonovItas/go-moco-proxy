package connpool

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/jackc/puddle/v2"
	"golang.org/x/sync/errgroup"
)

// Dialer is interface for connections establishment.
type Dialer interface {
	DialContext(ctx context.Context) (PoolConn, error)
}

// PoolConn is interface for pool connection.
type PoolConn interface {
	io.ReadWriteCloser
}

const (
	// defaultMinConns is default minimum count of connections.
	defaultMinConns int32 = 0

	// defaultMaxConn is default maximum count of connections.
	defaultMaxConns int32 = 4

	// defaultConnectTimeout is default timeout to establish connection.
	defaultConnectTimeout time.Duration = 5 * time.Second

	// defaultConnTimeout is default connection timeout.
	defaultConnTimeout time.Duration = 15 * time.Minute

	// defaultConnIdleTimeout is default idle connection timeout.
	defaultConnIdleTimeout time.Duration = 5 * time.Minute

	// defaultHealthCheckPeriod is default health check timeout.
	defaultHealthCheckPeriod time.Duration = 60 * time.Second

	// timeToDestroy is timeout to destoy resource in a pool.
	// That timeout is needed because puddle.Pool destroy resource concurrently.
	timeToDestroy time.Duration = 500 * time.Millisecond
)

// ConfigPool is config for connection pool.
type PoolConfig struct {
	// MinConns is count of minimum connections that should be in a pool.
	// If count of connections is lower than minConns, new connections would be
	// established after healthcheck period.
	MinConns int32

	// MaxConns is maximum connections that should be in a pool.
	MaxConns int32

	// ConnectTimeout is timeout to establish new connection.
	ConnectTimeout time.Duration

	// ConnTimeout is timeout since creation of a connection after
	// which a connection will be automatically closed.
	ConnTimeout time.Duration

	// ConnIdleTimeout is timeout after which idle connection will be automatically closed.
	ConnIdleTimeout time.Duration

	// HealthCheckPeriod is period between health cheks.
	HealthCheckPeriod time.Duration

	// ConnDialer is dialer for connetions.
	ConnDialer Dialer
}

// DefaultConfig function returns pool config with default parameters.
// Copy method returns copy of the pool config.
func (pc *PoolConfig) Copy() *PoolConfig {
	c := new(PoolConfig)
	*c = *pc
	return c
}

// defaultConfig function returns config with pool defaults.
func defaultConfig() *PoolConfig {
	return &PoolConfig{
		MinConns:          defaultMinConns,
		MaxConns:          defaultMaxConns,
		ConnectTimeout:    defaultConnectTimeout,
		ConnTimeout:       defaultConnTimeout,
		ConnIdleTimeout:   defaultConnIdleTimeout,
		HealthCheckPeriod: defaultHealthCheckPeriod,
	}
}

// TODO: add timeout checks
// validateConfig method checks correctness of the config.
func (pc *PoolConfig) validateConfig() error {
	if pc.ConnDialer == nil {
		return fmt.Errorf("dialer cannot be a nil value")
	}

	if pc.MaxConns <= 0 || pc.MinConns < 0 || pc.MaxConns < pc.MinConns {
		return fmt.Errorf(
			"invalid connections settings: min conns %d, max conns %d",
			pc.MinConns,
			pc.MaxConns,
		)
	}

	return nil
}

// connResource is struct that describes connection resource.
type connResource struct {
	PoolConn
}

// Pool is used for connection pooling and using pgx/puddle.Pool like a base pool.
type Pool struct {
	p          *puddle.Pool[*connResource]
	config     *PoolConfig
	connDialer Dialer

	minConns          int32
	maxConns          int32
	connTimeout       time.Duration
	connIdleTimeout   time.Duration
	healthCheckPeriod time.Duration
	healthCheckch     chan struct{}

	close   sync.Once
	closech chan struct{}
}

// NewPool function returns new pool with default config and given dialer.
func NewPool(d Dialer) (pool *Pool, err error) {
	if d == nil {
		return nil, fmt.Errorf("dialer cannot be nil value")
	}

	cfg := defaultConfig()
	cfg.ConnDialer = d

	return NewPoolWithConfig(cfg)
}

// NewPoolWithConfig function returns new pool with given config.
func NewPoolWithConfig(config *PoolConfig) (pool *Pool, err error) {
	if err := config.validateConfig(); err != nil {
		return nil, err
	}

	dialer := config.ConnDialer
	p, err := puddle.NewPool(&puddle.Config[*connResource]{
		Constructor: func(ctx context.Context) (connRes *connResource, err error) {
			timeoutCtx, cancel := context.WithTimeout(ctx, config.ConnectTimeout)
			defer cancel()

			conn, err := dialer.DialContext(timeoutCtx)
			if err != nil {
				return nil, err
			}

			return &connResource{conn}, nil
		},
		Destructor: func(connRes *connResource) {
			connRes.Close()
		},
		MaxSize: config.MaxConns,
	})
	if err != nil {
		return nil, err
	}

	pool = &Pool{
		p:          p,
		config:     config,
		connDialer: dialer,

		minConns:          config.MinConns,
		maxConns:          config.MaxConns,
		connTimeout:       config.ConnTimeout,
		connIdleTimeout:   config.ConnIdleTimeout,
		healthCheckPeriod: config.HealthCheckPeriod,

		healthCheckch: make(chan struct{}, 1),
		closech:       make(chan struct{}),
	}

	go pool.backgroundCheckHealth()

	return pool, nil
}

// Acquire method acquires new connection from the pool and returns it.
func (p *Pool) Acquire(ctx context.Context) (*Conn, error) {
	res, err := p.p.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	return &Conn{
		p:   p,
		res: res,
	}, nil
}

// AcquireFunc method acquires new connection and call f func with that connection.
// After f function is done, return connection to the pool.
func (p *Pool) AcquireFunc(ctx context.Context, f func(c *Conn) error) error {
	res, err := p.p.Acquire(ctx)
	if err != nil {
		return err
	}
	defer res.Release()

	return f(&Conn{
		p:   p,
		res: res,
	})
}

// Close method closes pool and wait until all connection is closed.
func (p *Pool) Close() {
	p.close.Do(func() {
		close(p.closech)
		p.p.Close()
	})
}

// CloneConfig method returns copy of pool config.
func (p *Pool) CloneConfig() *PoolConfig {
	return p.config.Copy()
}

func (p *Pool) Stat() *Stat {
	return &Stat{s: p.p.Stat()}
}

// triggerHealthCheck triggers pool health check.
func (p *Pool) triggerHealthCheck() {
	go func() {
		time.Sleep(timeToDestroy)
		select {
		case p.healthCheckch <- struct{}{}:
		default:
		}
	}()
}

// backgroundCheckHealth method starts health check every
// health check period and when it's triggered by health check channel.
func (p *Pool) backgroundCheckHealth() {
	ticker := time.NewTicker(p.healthCheckPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-p.closech:
			return
		case <-p.healthCheckch:
			p.checkHealth()
		case <-ticker.C:
			p.checkHealth()
		}
	}
}

// checkHealth method checks count of idle, destroyed and simple connections.
func (p *Pool) checkHealth() {
	for {
		if err := p.checkMinConns(); err != nil {
			break
		}

		if ok := p.checkConnsHealth(); ok {
			return
		}

		select {
		case <-p.closech:
		case <-time.After(timeToDestroy):
		}
	}
}

// checkConnsHealth method checks connections timeouts and if it's expired destroy resource
// and count of connections in pool is greater that minimum pool's connections.
func (p *Pool) checkConnsHealth() bool {
	healthy := true
	totalCounts := p.Stat().TotalConns()
	resources := p.p.AcquireAllIdle()
	for _, res := range resources {
		if p.isExpired(res) && totalCounts > p.minConns {
			healthy = false
			res.Destroy()
		} else if p.isIdleExpired(res) && totalCounts > p.minConns {
			healthy = false
			res.Destroy()
		} else {
			res.ReleaseUnused()
		}
	}

	return healthy
}

// checkMinConns method creates new connections if count of connections lower that minimum.
func (p *Pool) checkMinConns() error {
	toCreate := p.minConns - p.Stat().TotalConns()
	if toCreate > 0 {
		return p.createIdleConns(context.Background(), toCreate)
	}

	return nil
}

// createIdleConns method creates new idle connections in the pool.
func (p *Pool) createIdleConns(parentCtx context.Context, totalCounts int32) error {
	g, ctx := errgroup.WithContext(parentCtx)
	for i := 0; i < int(totalCounts); i++ {
		g.Go(func() error {
			return p.p.CreateResource(ctx)
		})
	}

	return g.Wait()
}

// isIdleExpired method checks connection is expired or not by connIdleTimeout.
func (p *Pool) isIdleExpired(res *puddle.Resource[*connResource]) bool {
	return isTimeoutExpired(res.IdleDuration(), p.connIdleTimeout)
}

// isExpired method checks connection is expired or not by connTimeout.
func (p *Pool) isExpired(res *puddle.Resource[*connResource]) bool {
	return isTimeoutExpired(time.Since(res.CreationTime()), p.connTimeout)
}
