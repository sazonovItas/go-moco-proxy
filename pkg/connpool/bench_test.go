package connpool

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockBenchDialer struct {
	mock.Mock
}

func (td *mockBenchDialer) DialContext(ctx context.Context) (PoolConn, error) {
	args := td.Called(ctx)
	return args.Get(0).(PoolConn), args.Error(1)
}

type benchConn struct {
	PoolConn
}

func (tc *benchConn) Close() error {
	return nil
}

func BenchmarkAcquireAndRelease(b *testing.B) {
	mockDialer := new(mockBenchDialer)
	mockDialer.On("DialContext", mock.Anything).Return(new(benchConn), nil)

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

	mockDialer := new(mockBenchDialer)
	mockDialer.On("DialContext", mock.Anything).Return(new(benchConn), nil)
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
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
