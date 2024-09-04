package proxy

import (
	"context"
	"sync"
)

type Proxy struct {
	cancel context.CancelFunc
	sync.Mutex

	wg sync.WaitGroup
}

func NewProxy() *Proxy {
	return &Proxy{}
}

func (p *Proxy) ListenAndServe(ctx context.Context) {
	p.Lock()
	if p.cancel != nil {
		p.cancel()
		p.cancel = nil
	}

	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	p.Unlock()
}

func (p *Proxy) Shutdown() {
	p.Lock()
	if p.cancel != nil {
		p.cancel()
		p.cancel = nil
	}
	p.Unlock()
}
