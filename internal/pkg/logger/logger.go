package logger

import (
	"context"
	"log/slog"
	"os"
	"reflect"
)

type HandlerMiddleware struct {
	next slog.Handler
}

func NewHandlerMiddleware(next slog.Handler) *HandlerMiddleware {
	return &HandlerMiddleware{next: next}
}

func (h *HandlerMiddleware) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}

func (h *HandlerMiddleware) Handle(ctx context.Context, rec slog.Record) error {
	if ctxValue, ok := ctx.Value(key).(loggerCtx); ok {
		val := reflect.ValueOf(ctxValue)
		typ := reflect.TypeOf(ctxValue)

		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldType := typ.Field(i)

			// Пропускаем неэкспортированные поля
			if !fieldType.IsExported() {
				continue
			}

			if !field.IsZero() {
				rec.Add(fieldType.Name, field.Interface())
			}
		}
	}

	return h.next.Handle(ctx, rec)
}

func (h *HandlerMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddleware{next: h.next.WithAttrs(attrs)} // не забыть обернуть, но осторожно
}

func (h *HandlerMiddleware) WithGroup(name string) slog.Handler {
	return &HandlerMiddleware{next: h.next.WithGroup(name)} // не забыть обернуть, но осторожно
}

func InitLogger() {
	handler := slog.Handler(NewPrettyTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	handler = NewHandlerMiddleware(handler)
	slog.SetDefault(slog.New(handler))
}
