package metrics

import (
	"testing"
	"time"

	"github.com/sazonovItas/go-moco-proxy/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name         string
		cfg          config.MetricConfig
		expectedAddr string
	}{{
		name: "valid metric's configuration",
		cfg: config.MetricConfig{
			Address: "localhost:8080",
		},
		expectedAddr: "localhost:8080",
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := New(tc.cfg)
			require.NoError(t, err, "should not be error create metric's server")
			require.Equal(
				t,
				tc.expectedAddr,
				m.server.Addr,
				"expected %s, got %s should equal addresses",
				tc.expectedAddr,
				m.server.Addr,
			)
		})
	}
}

func TestRunAndShutdown(t *testing.T) {
	const timeoutToRun time.Duration = 250 * time.Millisecond

	testCases := []struct {
		name            string
		cfg             config.MetricConfig
		wantRunErr      bool
		wantShutdownErr bool
	}{
		{
			name: "valid metric's server run",
			cfg: config.MetricConfig{
				Address: "127.0.0.1:8080",
			},
			wantRunErr:      false,
			wantShutdownErr: false,
		},
		{
			name: "invalid metric's configuration",
			cfg: config.MetricConfig{
				Address: "_kHiohelh_:8080",
			},
			wantRunErr:      true,
			wantShutdownErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := New(tc.cfg)
			require.NoError(t, err, "should not be error create metric's server")

			errRunCh := make(chan error)
			defer close(errRunCh)

			go func() {
				errRunCh <- m.Run()
			}()

			<-time.After(timeoutToRun)
			err = m.Shutdown()
			if tc.wantShutdownErr {
				require.Error(t, err, "no error occurred while shutdown server")
				return
			}
			require.NoError(t, err, "error occured while running server: %v", err)

			err = <-errRunCh
			if tc.wantRunErr {
				require.Error(t, err, "no error occurred while running server")
				return
			}
			require.NoError(t, err, "error occurred while running server: %v", err)
		})
	}
}

func TestShutdownOnNotRunningServer(t *testing.T) {
	m, err := New(config.MetricConfig{})
	err = m.Shutdown()
	require.NoError(t, err, "shouldn't be error to shutdown not running serrver")
}
