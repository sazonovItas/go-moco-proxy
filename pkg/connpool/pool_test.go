package connpool

import (
	"context"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testDialer struct {
	mock.Mock
}

func (td *testDialer) DialContext(ctx context.Context) (PoolConn, error) {
	return nil, nil
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
				ConnDialer:        &testDialer{},
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
