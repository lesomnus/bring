package bringer

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/lesomnus/bring/log"
	"github.com/lesomnus/bring/thing"
)

type httpBringer struct{}

func HttpBringer(opts ...Option) Bringer {
	return &httpBringer{}
}

func (b *httpBringer) Bring(ctx context.Context, t thing.Thing) (io.ReadCloser, error) {
	l := log.From(ctx).With(name("http"))

	// TODO: check ETag, Cache-Control, of Last-Modified header.
	// res, err := http.Head(t.Url.String())
	// if err != nil {
	// 	return nil, fmt.Errorf("request head: %w", err)
	// }

	// TODO: do not buffer the response in the memory according to the config.
	// f, err := os.CreateTemp("", "bring-")
	// if err != nil {
	// 	return nil, fmt.Errorf("create temp file: %w", err)
	// }
	// defer f.Close()

	l.Info("request GET")
	res, err := http.Get(t.Url.String())
	if err != nil {
		return nil, fmt.Errorf("request get: %w", err)
	}
	l.Info("response", slog.Int("status", res.StatusCode))
	if l.Enabled(ctx, slog.LevelDebug) {
		l.Debug("response", slog.String("header", fmt.Sprintf("%v", res.Header)))
	}

	return res.Body, nil
}
