package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var loggerConfig = zap.NewProductionConfig()

type (
	LoggerOption interface {
		Apply(cfg *zap.Config)
	}

	LoggerOptionFunc func(cfg *zap.Config)
)

func (lf LoggerOptionFunc) Apply(cfg *zap.Config) {
	lf(cfg)
}

func WithLevel(level zap.AtomicLevel) LoggerOptionFunc {
	return func(cfg *zap.Config) {
		cfg.Level = level
	}
}

func WithDevelopmentLogs(development bool) LoggerOptionFunc {
	return func(cfg *zap.Config) {
		cfg.Development = development
	}
}

// WithEncoding function returns new logger option.
// Valid values for encoding are "json" and "console".
func WithEncoding(encoding string) LoggerOptionFunc {
	return func(cfg *zap.Config) {
		cfg.Encoding = encoding
	}
}

func WithInitialFields(fields map[string]any) LoggerOptionFunc {
	return func(cfg *zap.Config) {
		cfg.InitialFields = fields
	}
}

func WithOutputPaths(paths []string) LoggerOptionFunc {
	return func(cfg *zap.Config) {
		cfg.OutputPaths = paths
	}
}

func WithErrorOutputPaths(paths []string) LoggerOptionFunc {
	return func(cfg *zap.Config) {
		cfg.ErrorOutputPaths = paths
	}
}

func WithEncoderConfig(encoder zapcore.EncoderConfig) LoggerOptionFunc {
	return func(cfg *zap.Config) {
		cfg.EncoderConfig = encoder
	}
}

// ConfigureLogger function configures logger config with options
// and replaces global zap logger with built one.
func ConfigureLogger(options ...LoggerOption) error {
	for _, option := range options {
		option.Apply(&loggerConfig)
	}

	if l, err := loggerConfig.Build(zap.AddCallerSkip(1)); err != nil {
		return err
	} else {
		zap.ReplaceGlobals(l)
	}

	return nil
}

// CreateLogger function returns zap global logger.
func CreateLogger() *zap.Logger {
	return zap.L()
}

// ParseLevel functions parses level from string.
func ParseLevel(level string) zap.AtomicLevel {
	switch strings.ToLower(level) {
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		return zap.NewAtomicLevelAt(zap.FatalLevel)
	case "panic":
		return zap.NewAtomicLevelAt(zap.PanicLevel)
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
}
