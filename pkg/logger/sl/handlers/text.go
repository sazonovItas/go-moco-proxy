package slhandlers

import (
	"context"
	"encoding/json"
	"io"
	stdLog "log"
	"log/slog"

	"github.com/fatih/color"
)

// TextHandlerOptions struct of console handler options.
type TextHandlerOptions struct {
	SlogOpts *slog.HandlerOptions
}

// TextHandler struct o console handler.
type TextHandler struct {
	slog.Handler
	l     *stdLog.Logger
	attrs []slog.Attr
}

// NewTextHandler method returns new console handler to log at given out writer.
func (opts TextHandlerOptions) NewTextHandler(out io.Writer) *TextHandler {
	return &TextHandler{
		Handler: slog.NewTextHandler(out, opts.SlogOpts),
		l:       stdLog.New(out, "", 0),
	}
}

// Enabled method is implementation of slog handler interface.
func (h *TextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

// Handler method is implementation of slog handler interface.
func (h *TextHandler) Handle(ctx context.Context, r slog.Record) error {
	level := "[" + r.Level.String() + "]"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	var b []byte
	if len(fields) > 0 {
		for key, value := range fields {
			text, err := json.Marshal(value)
			if err != nil {
				return err
			}

			b = append(b, []byte(key+": "+string(text)+", ")...)
		}
	}

	time := color.MagentaString(r.Time.Format("[02/01 2006 15:04:05.000]"))
	msg := color.CyanString(r.Message)

	h.l.Println(time, level, msg, string(b))

	return nil
}

// WithAttrs method is implementation of slog handler interface.
func (h *TextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TextHandler{
		Handler: h.Handler.WithAttrs(attrs),
		l:       h.l,
		attrs:   attrs,
	}
}

// WithGroup method is implementation of slog handler interface.
func (h *TextHandler) WithGroup(name string) slog.Handler {
	return &TextHandler{
		Handler: h.Handler.WithGroup(name),
		l:       h.l,
	}
}
