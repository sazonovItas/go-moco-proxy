package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/sazonovItas/go-moco-proxy/pkg/logger/sl"
)

// Level is type for logger level.
type Level uint8

const (
	LevelDebug Level = iota // debug-level
	LevelInfo               // info-level
	LevelWarn               // warn-level
	LevelError              // error-level
	LevelFatal              // fatal-level
	LevelPanic              // panic-level
)

var (
	log      *slog.Logger
	logLevel slog.Level = slog.LevelInfo

	logOnce sync.Once
)

// SetLoggerLevel sets logger level for global logger.
func SetLoggerLevel(level Level) {
	logLevel = convertLoggerLevel(level)
}

func Log(msg string, level Level, args ...any) {
	getLogger().Log(context.Background(), convertLoggerLevel(level), msg, args...)
}

func Debug(msg string, args ...any) {
	getLogger().Debug(msg, args...)
}

func Info(msg string, args ...any) {
	getLogger().Info(msg, args...)
}

func Warn(msg string, args ...any) {
	getLogger().Warn(msg, args...)
}

func Error(msg string, args ...any) {
	getLogger().Error(msg, args...)
}

func Debugf(msg string, args ...any) {
	getLogger().Debug(fmt.Sprintf(msg, args...))
}

func Infof(msg string, args ...any) {
	getLogger().Info(fmt.Sprintf(msg, args...))
}

func Warnf(msg string, args ...any) {
	getLogger().Warn(fmt.Sprintf(msg, args...))
}

func Errorf(msg string, args ...any) {
	getLogger().Error(fmt.Sprintf(msg, args...))
}

// getLogger returns singleton logger.
func getLogger() *slog.Logger {
	logOnce.Do(func() {
		log = sl.NewLogger(logLevel, os.Stdout)
	})

	return log
}

// convertLoggerLevel converts logger level to slog.Level.
func convertLoggerLevel(level Level) slog.Level {
	if level > LevelError || level < LevelDebug {
		return slog.LevelDebug
	}

	return [...]slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError}[level]
}
