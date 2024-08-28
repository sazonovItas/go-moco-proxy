package connpool

import (
	"context"
	"math/rand"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockDialer struct {
	mock.Mock
}

func (td *mockDialer) DialContext(ctx context.Context) (PoolConn, error) {
	args := td.Called(ctx)
	return args.Get(0).(PoolConn), args.Error(1)
}

type mockPoolConn struct {
	mock.Mock
}

func (pc *mockPoolConn) Read(p []byte) (n int, err error) {
	return len(p), nil
}

func (pc *mockPoolConn) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (pc *mockPoolConn) Close() error {
	args := pc.Called()
	return args.Error(0)
}

func TestNewPoolWithConfig(t *testing.T) {
	t.Parallel()

	mockPoolConn := new(mockPoolConn)
	mockPoolConn.On("Close").Return(nil)

	mockDialer := new(mockDialer)
	mockDialer.On("DialContext", mock.Anything).Return(nil, nil)

	testCases := []struct {
		name    string
		cfg     *PoolConfig
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &PoolConfig{
				MinConns:          defaultMinConns,
				MaxConns:          defaultMaxConns,
				ConnectTimeout:    defaultConnectTimeout,
				ConnTimeout:       defaultConnTimeout,
				ConnIdleTimeout:   defaultConnIdleTimeout,
				HealthCheckPeriod: defaultHealthCheckPeriod,
				ConnDialer:        mockDialer,
			},
			wantErr: false,
		},
		{
			name:    "config dialer nil value",
			cfg:     defaultConfig(),
			wantErr: true,
		},
		{
			name:    "empty config",
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
	const N = 1000

	mockPoolConn := new(mockPoolConn)
	mockPoolConn.On("Close").Return(nil)

	mockDialer := new(mockDialer)
	mockDialer.On("DialContext", mock.Anything).Return(mockPoolConn, nil)

	pool, err := NewPool(mockDialer)
	require.NoError(t, err, "failed create pool: %v", err)
	defer pool.Close()

	for i := 0; i < N; i++ {
		buf := make([]byte, rand.Intn(N))

		conn, err := pool.Acquire(context.Background())
		if err != nil {
			t.Fatalf("failed acquire connection: %v", err)
		}

		n, _ := conn.Conn().Write(buf)
		if n != len(buf) {
			t.Fatalf("failed write to connection: want %d, got %d", len(buf), n)
		}
		conn.Release()

		conn, err = pool.Acquire(context.Background())
		if err != nil {
			t.Fatalf("failed acquire connection: %v", err)
		}

		c := conn.Hijack()
		require.NotNil(t, c, "failed hijack connection")

		conn.Release()
	}
}

func TestPoolAcquireFunc(t *testing.T) {
	const N = 1000

	mockPoolConn := new(mockPoolConn)
	mockPoolConn.On("Close").Return(nil)

	mockDialer := new(mockDialer)
	mockDialer.On("DialContext", mock.Anything).Return(mockPoolConn, nil)

	pool, err := NewPool(mockDialer)
	require.NoError(t, err, "failed create pool: %v", err)
	defer pool.Close()

	for i := 0; i < N; i++ {
		buf := make([]byte, rand.Intn(N))

		err := pool.AcquireFunc(context.Background(), func(c *Conn) error {
			n, _ := c.Conn().Write(buf)
			if n != len(buf) {
				t.Fatalf("failed write to connection: want %d, got %d", len(buf), n)
			}

			return nil
		})
		if err != nil {
			t.Fatalf("failed acquire connection: %v", err)
		}
	}
}

func TestPoolClose(t *testing.T) {
	t.Parallel()

	mockPoolConn := new(mockPoolConn)
	mockPoolConn.On("Close").Return(nil)

	mockDialer := new(mockDialer)
	mockDialer.On("DialContext", mock.Anything).Return(mockPoolConn, nil)

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
			name: "valid config",
			cfg: &PoolConfig{
				MinConns:          defaultMinConns,
				MaxConns:          defaultMaxConns,
				ConnectTimeout:    defaultConnectTimeout,
				ConnTimeout:       defaultConnTimeout,
				ConnIdleTimeout:   defaultConnIdleTimeout,
				HealthCheckPeriod: defaultHealthCheckPeriod,
				ConnDialer:        &mockDialer{},
			},
			wantErr: false,
		},
		{
			name:    "config dialer nil value",
			cfg:     defaultConfig(),
			wantErr: true,
		},
		{
			name:    "empty config",
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
