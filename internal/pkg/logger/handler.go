package logger

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
)

const (
	timeFormat = "[2006-01-02 15:04:05.000]"
)

type PrettyTextHandler struct {
	next   slog.Handler
	buffer bytes.Buffer
}

func NewPrettyTextHandler(w io.Writer, opts *slog.HandlerOptions) *PrettyTextHandler {
	return &PrettyTextHandler{next: slog.NewTextHandler(w, opts)}
}

func (h *PrettyTextHandler) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}

func (h *PrettyTextHandler) Handle(ctx context.Context, rec slog.Record) error {
	h.buffer.Reset()

	h.buffer.WriteString(colorize(colorCodeLightGray, rec.Time.Format(timeFormat)+" "))

	var level = rec.Level.String() + ": "
	switch rec.Level {
	case slog.LevelDebug:
		level = colorize(colorCodeDarkGray, level)
	case slog.LevelInfo:
		level = colorize(colorCodeCyan, level)
	case slog.LevelWarn:
		level = colorize(colorCodeLightYellow, level)
	case slog.LevelError:
		level = colorize(colorCodeLightRed, level)
	}

	h.buffer.WriteString(level)
	h.buffer.WriteString(rec.Message)

	rec.Attrs(func(attr slog.Attr) bool {
		h.buffer.WriteString(" " + colorize(colorCodeBlue, attr.String()))
		return true
	})

	fmt.Println(h.buffer.String())

	return nil
}

func (h *PrettyTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyTextHandler{next: h.next.WithAttrs(attrs)}
}

func (h *PrettyTextHandler) WithGroup(name string) slog.Handler {
	return &PrettyTextHandler{next: h.next.WithGroup(name)}
}
