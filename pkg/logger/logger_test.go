package logger

import (
	"log/slog"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func Test_getLogger(t *testing.T) {
	require.Equal(
		t,
		unsafe.Pointer(getLogger()),
		unsafe.Pointer(getLogger()),
		"should be equal pointer addresses",
	)
}

func TestSetLoggerLevel(t *testing.T) {
	testCases := []struct {
		name  string
		level Level
		want  slog.Level
	}{
		{
			name:  "debug logger level",
			level: LevelDebug,
			want:  slog.LevelDebug,
		},
		{
			name:  "info logger level",
			level: LevelInfo,
			want:  slog.LevelInfo,
		},
		{
			name:  "warning logger level",
			level: LevelWarn,
			want:  slog.LevelWarn,
		},
		{
			name:  "error logger level",
			level: LevelError,
			want:  slog.LevelError,
		},
		{
			name:  "unknow logger level",
			level: 20,
			want:  slog.LevelDebug,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			SetLoggerLevel(tc.level)
			if got := logLevel; got != tc.want {
				t.Fatalf("got %s, want %s", got.String(), tc.want.String())
			}
		})
	}
}

func Test_convertLoggerLevel(t *testing.T) {
	testCases := []struct {
		name  string
		level Level
		want  slog.Level
	}{
		{
			name:  "debug logger level",
			level: LevelDebug,
			want:  slog.LevelDebug,
		},
		{
			name:  "info logger level",
			level: LevelInfo,
			want:  slog.LevelInfo,
		},
		{
			name:  "warning logger level",
			level: LevelWarn,
			want:  slog.LevelWarn,
		},
		{
			name:  "error logger level",
			level: LevelError,
			want:  slog.LevelError,
		},
		{
			name:  "unknow logger level",
			level: 20,
			want:  slog.LevelDebug,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := convertLoggerLevel(tc.level); got != tc.want {
				t.Fatalf("got %s, want %s", got.String(), tc.want.String())
			}
		})
	}
}
