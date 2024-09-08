package logger

import (
	"go.uber.org/zap"
)

func Sync() error {
	return getLogger().Sync()
}

func Named(name string) *zap.Logger {
	return getLogger().Named(name)
}

func With(fields ...zap.Field) *zap.Logger {
	return getLogger().With(fields...)
}

func Debug(msg string, fields ...zap.Field) {
	getLogger().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	getLogger().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	getLogger().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	getLogger().Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	getLogger().Fatal(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	getLogger().Panic(msg, fields...)
}

func Sugar() *zap.SugaredLogger {
	return getLogger().Sugar()
}

func getLogger() *zap.Logger {
	return CreateLogger()
}
