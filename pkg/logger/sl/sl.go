package sl

import (
	"io"
	"log/slog"

	slhandlers "github.com/sazonovItas/go-moco-proxy/pkg/logger/sl/handlers"
)

// NewLogger functions creates new logger depence on logger type json and text.
func NewLogger(level slog.Level, out io.Writer) (log *slog.Logger) {
	return NewConsoleLogger(level, out)
}

// NewTextLogger function creates text handler for slog logger.
func NewConsoleLogger(level slog.Level, out io.Writer) *slog.Logger {
	opts := slhandlers.TextHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: level,
		},
	}

	handler := opts.NewTextHandler(out)
	return slog.New(handler)
}

// Err function returns slog attribute from error.
func Err(err error) slog.Attr {
	if err == nil {
		slog.String("error", "nil")
	}

	return slog.String("error", err.Error())
}
