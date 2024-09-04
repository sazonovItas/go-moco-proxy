package connpool

import (
	"context"
	"net"
	"testing"
	"unsafe"

	mocknet "github.com/sazonovItas/go-moco-proxy/mocks/net"
	mockconnpool "github.com/sazonovItas/go-moco-proxy/pkg/connpool/mock/connpool"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewPoolWithConfig(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		cfg     *PoolConfig
		wantErr bool
	}{
		{
			name: "Test valid config",
			cfg: &PoolConfig{
				MinConns:          defaultMinConns,
				MaxConns:          defaultMaxConns,
				ConnectTimeout:    defaultConnectTimeout,
				ConnTimeout:       defaultConnTimeout,
				ConnIdleTimeout:   defaultConnIdleTimeout,
				HealthCheckPeriod: defaultHealthCheckPeriod,
				ConnDialer:        &mockconnpool.MockDialer{},
			},
			wantErr: false,
		},
		{
			name:    "Test config dialer nil value",
			cfg:     defaultConfig(),
			wantErr: true,
		},
		{
			name:    "Test empty config",
			cfg:     &PoolConfig{},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pool, err := NewPoolWithConfig(tc.cfg)
			if err != nil {
				if !tc.wantErr {
					t.Fatalf("failed to create pool: %v", err)
				}
			} else {
				if tc.wantErr {
					t.Fatalf("expected error to create pool")
				}

				pool.Close()
			}
		})
	}
}

func TestPoolAcquireAndHijack(t *testing.T) {
	const N = 10

	mockDialer := mockconnpool.NewMockDialer(t)
	mockDialer.EXPECT().
		DialContext(mock.Anything).
		RunAndReturn(func(ctx context.Context) (net.Conn, error) {
			mockConn := mocknet.NewMockConn(t)
			mockConn.EXPECT().Close().Once().Return(nil)
			return mockConn, nil
		})

	pool, err := NewPool(mockDialer)
	require.NoError(t, err, "failed create pool: %v", err)
	defer pool.Close()

	for i := 0; i < N; i++ {
		conn, err := pool.Acquire(context.Background())
		if err != nil {
			t.Fatalf("failed acquire connection: %v", err)
		}
		conn.Release()

		conn, err = pool.Acquire(context.Background())
		if err != nil {
			t.Fatalf("failed acquire connection: %v", err)
		}

		c := conn.Hijack()
		defer c.Close()
		require.NotNil(t, c, "failed hijack connection")

		conn.Release()
	}
}

func TestPoolAcquireFunc(t *testing.T) {
	const N = 10

	mockDialer := mockconnpool.NewMockDialer(t)
	mockDialer.EXPECT().
		DialContext(mock.Anything).
		RunAndReturn(func(ctx context.Context) (net.Conn, error) {
			mockConn := mocknet.NewMockConn(t)
			mockConn.EXPECT().Close().Once().Return(nil)
			return mockConn, nil
		})

	pool, err := NewPool(mockDialer)
	require.NoError(t, err, "failed create pool: %v", err)
	defer pool.Close()

	for i := 0; i < N; i++ {
		err := pool.AcquireFunc(context.Background(), func(c *Conn) error {
			return nil
		})
		if err != nil {
			t.Fatalf("failed acquire connection: %v", err)
		}
	}
}

func TestPoolClose(t *testing.T) {
	t.Parallel()

	mockDialer := mockconnpool.NewMockDialer(t)
	mockDialer.EXPECT().
		DialContext(mock.Anything).
		RunAndReturn(func(ctx context.Context) (net.Conn, error) {
			mockConn := mocknet.NewMockConn(t)
			mockConn.EXPECT().Close().Once().Return(nil)
			return mockConn, nil
		}).Maybe()

	pool, err := NewPool(mockDialer)
	require.NoError(t, err, "failed create pool: %w", err)

	pool.Close()

	_, ok := <-pool.closech
	require.False(t, ok, "close chanel should be closed")
}

func TestCopyConfig(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		want *PoolConfig
	}{
		{
			name: "copy default config",
			want: defaultConfig(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.want.Copy()
			require.Equal(t, tc.want, got, "should be equal configs")
			require.NotEqual(
				t,
				unsafe.Pointer(tc.want),
				unsafe.Pointer(got),
				"should be difference instances",
			)
		})
	}
}

func Test_validateConfig(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		cfg     *PoolConfig
		wantErr bool
	}{
		{
			name: "Test valid config",
			cfg: &PoolConfig{
				MinConns:          defaultMinConns,
				MaxConns:          defaultMaxConns,
				ConnectTimeout:    defaultConnectTimeout,
				ConnTimeout:       defaultConnTimeout,
				ConnIdleTimeout:   defaultConnIdleTimeout,
				HealthCheckPeriod: defaultHealthCheckPeriod,
				ConnDialer:        &mockconnpool.MockDialer{},
			},
			wantErr: false,
		},
		{
			name:    "Test config dialer nil value",
			cfg:     defaultConfig(),
			wantErr: true,
		},
		{
			name:    "Test empty config",
			cfg:     &PoolConfig{},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cfg.validateConfig()
			if tc.wantErr {
				require.Error(t, err, "should be not valid config: %v", tc.cfg)
			} else {
				require.NoError(t, err, "should valid config: %v", tc.cfg)
			}
		})
	}
}
