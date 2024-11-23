package log

import (
	"context"
	"io"
	"log/slog"
)

type logCtxKey struct{}

var Discard = slog.New(slog.NewTextHandler(io.Discard, nil))

func From(ctx context.Context) *slog.Logger {
	v, ok := ctx.Value(logCtxKey{}).(*slog.Logger)
	if !ok {
		return Discard
	}

	return v
}

func Into(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, logCtxKey{}, logger)
}
