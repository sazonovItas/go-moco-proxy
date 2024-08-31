package logger

import "go.uber.org/zap"

type Logger interface {
	Sync() error
	Named(name string) Logger
	With(fields ...zap.Field) Logger
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
}
