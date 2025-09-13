package logger

import (
	"context"
	"errors"
)

type errorWithLogCtx struct {
	next error
	ctx  loggerCtx
}

func (e *errorWithLogCtx) Unwrap() error {
	return e.next
}

func (e *errorWithLogCtx) Error() string {
	return e.next.Error()
}

func WrapError(ctx context.Context, err error) error {
	var loggerErr *errorWithLogCtx
	if errors.As(err, &loggerErr) {
		return loggerErr
	}

	c := loggerCtx{}
	if x, ok := ctx.Value(key).(loggerCtx); ok {
		c = x
	}
	return &errorWithLogCtx{
		next: err,
		ctx:  c,
	}
}

func CtxFromError(ctx context.Context, err error) context.Context {
	var loggerErr *errorWithLogCtx

	if errors.As(err, &loggerErr) {
		return context.WithValue(ctx, key, loggerErr.ctx)
	}
	return ctx
}
