package logger

import "go.uber.org/zap"

type logger struct {
	ulogger *zap.Logger
}

var _ Logger = (*logger)(nil)

func NewLogger(l *zap.Logger) *logger {
	return &logger{ulogger: l}
}

func (l *logger) Named(name string) Logger {
	return NewLogger(l.ulogger.Named(name))
}

func (l *logger) With(fields ...zap.Field) Logger {
	return NewLogger(l.ulogger.With(fields...))
}

func (l *logger) Debug(msg string, fields ...zap.Field) {
	l.ulogger.Debug(msg, fields...)
}

func (l *logger) Info(msg string, fields ...zap.Field) {
	l.ulogger.Info(msg, fields...)
}

func (l *logger) Warn(msg string, fields ...zap.Field) {
	l.ulogger.Warn(msg, fields...)
}

func (l *logger) Error(msg string, fields ...zap.Field) {
	l.ulogger.Error(msg, fields...)
}

func (l *logger) Fatal(msg string, fields ...zap.Field) {
	l.ulogger.Fatal(msg, fields...)
}

func (l *logger) Panic(msg string, fields ...zap.Field) {
	l.ulogger.Panic(msg, fields...)
}
