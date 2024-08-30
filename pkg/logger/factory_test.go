package logger

import (
	"reflect"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestParseLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		level string
		want  zap.AtomicLevel
	}{
		{
			name:  "parse DEBUG level",
			level: "DEBUG",
			want:  zap.NewAtomicLevelAt(zap.DebugLevel),
		},
		{
			name:  "parse DeBuG level",
			level: "DeBuG",
			want:  zap.NewAtomicLevelAt(zap.DebugLevel),
		},
		{
			name:  "parse INFO level",
			level: "INFO",
			want:  zap.NewAtomicLevelAt(zap.InfoLevel),
		},
		{
			name:  "parse info level",
			level: "info",
			want:  zap.NewAtomicLevelAt(zap.InfoLevel),
		},
		{
			name:  "parse WARN level",
			level: "WARN",
			want:  zap.NewAtomicLevelAt(zap.WarnLevel),
		},
		{
			name:  "parse ERROR level",
			level: "ERROR",
			want:  zap.NewAtomicLevelAt(zap.ErrorLevel),
		},
		{
			name:  "parse Fatal level",
			level: "Fatal",
			want:  zap.NewAtomicLevelAt(zap.FatalLevel),
		},
		{
			name:  "parse panic level",
			level: "panic",
			want:  zap.NewAtomicLevelAt(zap.PanicLevel),
		},
		{
			name:  "parse unknown level",
			level: "unknown",
			want:  zap.NewAtomicLevelAt(zap.InfoLevel),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseLevel(tt.level); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigureLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		opts    []LoggerOption
		wantErr bool
	}{
		{
			name: "configure defaults",
		},
		{
			name: "configure with level",
			opts: []LoggerOption{WithLevel(zap.NewAtomicLevelAt(zap.FatalLevel))},
		},
		{
			name: "configure with development",
			opts: []LoggerOption{WithDevelopmentLogs(true)},
		},
		{
			name: "configure with encoding",
			opts: []LoggerOption{WithEncoding("console")},
		},
		{
			name: "configure with initial fields",
			opts: []LoggerOption{WithInitialFields(map[string]any{
				"name": []string{"1", "2", "3"},
				"time": time.Now(),
			})},
		},
		{
			name: "configure with ouput paths",
			opts: []LoggerOption{WithOutputPaths([]string{"stdout"})},
		},
		{
			name: "configure with error ouput paths",
			opts: []LoggerOption{WithErrorOutputPaths([]string{"stderr"})},
		},
		{
			name: "configure with error ouput paths",
			opts: []LoggerOption{WithEncoderConfig(zap.NewDevelopmentConfig().EncoderConfig)},
		},
		{
			name: "configure with several options",
			opts: []LoggerOption{
				WithLevel(zap.NewAtomicLevelAt(zap.ErrorLevel)),
				WithDevelopmentLogs(true),
				WithEncoding("json"),
			},
		},
		{
			name: "configure with error encoding",
			opts: []LoggerOption{
				WithEncoding("text"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ConfigureLogger(tt.opts...); (err != nil) != tt.wantErr {
				t.Fatalf("ConfigureLogger() error = %v", err)
			}
		})
	}
}
