package dialer

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"sync"
)

var ErrNoTargetAvaliable = errors.New("cannot dial with targets")

type BaseDialer interface {
	DialContext(ctx context.Context, network string, address string) (conn net.Conn, err error)
}

type Target struct {
	Host string
	Port string
}

// TODO: add errors dial errors to load balancer or
// move this dialer to proxy (because that core logic of proxy)
type LoadBalanceDialer struct {
	mu         sync.Mutex
	targets    []Target
	currTarget int

	network string
	dialer  BaseDialer
}

// NewLoadBalanceDialer returns new load balancer dialer.
func NewLoadBalanceDialer(network string, targets []Target) *LoadBalanceDialer {
	d := &LoadBalanceDialer{
		targets: targets,
		network: network,
	}

	if d.dialer == nil {
		d.dialer = &net.Dialer{}
	}

	return d
}

// DialContext method returns new connection with target.
func (lb *LoadBalanceDialer) DialContext(ctx context.Context) (net.Conn, error) {
	maxTries := len(lb.targets)
	for i := 0; i < maxTries; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()

		// chose next target and try dial with it
		default:
			target := lb.nextTarget()

			conn, err := lb.dialer.DialContext(
				ctx,
				lb.network,
				net.JoinHostPort(target.Host, target.Port),
			)
			if err == nil {
				return conn, nil
			}
		}
	}

	return nil, ErrNoTargetAvaliable
}

// nextTarget method choses next target for dialing.
func (lb *LoadBalanceDialer) nextTarget() Target {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	targets := lb.targets
	if lb.currTarget == len(targets) {
		lb.currTarget = 0
		rand.Shuffle(len(targets), func(i, j int) {
			targets[i], targets[j] = targets[j], targets[i]
		})
	}

	return targets[lb.currTarget]
}
