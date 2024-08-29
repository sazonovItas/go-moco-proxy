package connpool

import (
	"context"
	"net"
	"testing"

	mockconnpool "github.com/sazonovItas/go-moco-proxy/mocks/connpool"
	mocknet "github.com/sazonovItas/go-moco-proxy/mocks/net"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func BenchmarkAcquireAndRelease(b *testing.B) {
	mockDialer := mockconnpool.NewMockDialer(b)
	mockDialer.EXPECT().
		DialContext(mock.Anything).
		RunAndReturn(func(ctx context.Context) (net.Conn, error) {
			mockConn := mocknet.NewMockConn(b)
			mockConn.EXPECT().Close().Once().Return(nil)
			return mockConn, nil
		})

	pool, err := NewPool(mockDialer)
	require.NoError(b, err)
	defer pool.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conn, err := pool.Acquire(context.Background())
		if err != nil {
			b.Fatalf("error occurred while acquire: %v", err)
		}

		conn.Release()
	}

	b.ReportAllocs()
}

func BenchmarkAcquireAndReleaseParallel(b *testing.B) {
	testCases := []struct {
		name     string
		maxConns int32
	}{
		{
			name:     "pool with 4 maximum connections",
			maxConns: 4,
		},
		{
			name:     "pool with 8 maximum connections",
			maxConns: 8,
		},
		{
			name:     "pool with 16 maximum connections",
			maxConns: 16,
		},
		{
			name:     "pool with 32 maximum connections",
			maxConns: 32,
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			mockDialer := mockconnpool.NewMockDialer(b)
			mockDialer.EXPECT().
				DialContext(mock.Anything).
				RunAndReturn(func(ctx context.Context) (net.Conn, error) {
					mockConn := mocknet.NewMockConn(b)
					mockConn.EXPECT().Close().Once().Return(nil)
					return mockConn, nil
				})
			pool, err := NewPoolWithConfig(&PoolConfig{
				MinConns:          defaultMinConns,
				MaxConns:          tc.maxConns,
				ConnectTimeout:    defaultConnTimeout,
				ConnTimeout:       defaultConnTimeout,
				ConnIdleTimeout:   defaultConnIdleTimeout,
				HealthCheckPeriod: defaultHealthCheckPeriod,
				ConnDialer:        mockDialer,
			})
			require.NoError(b, err)
			defer pool.Close()

			b.ResetTimer()
			b.RunParallel(func(p *testing.PB) {
				for p.Next() {
					conn, err := pool.Acquire(context.Background())
					if err != nil {
						b.Fatalf("error occurred while acquire: %v", err)
					}

					conn.Release()
				}
			})

			b.ReportAllocs()
		})
	}
}
