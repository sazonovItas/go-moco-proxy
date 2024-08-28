package connpool

import (
	"context"
	"io"
)

// Dialer is interface for connections establishment.
type Dialer interface {
	DialContext(ctx context.Context) (PoolConn, error)
}

// PoolConn is interface for pool connection.
type PoolConn interface {
	io.ReadWriteCloser
}
