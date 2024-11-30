package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/lesomnus/bring/bringer"
	"github.com/lesomnus/bring/log"
	"github.com/lesomnus/bring/thing"
)

type br struct {
	client http.Client
}

func (b *br) apply(opts []bringer.Option) {
	for _, opt := range opts {
		switch o := opt.(type) {
		case *transportOption:
			b.client.Transport = o.Value
		}
	}
}

func (b *br) with(opts []bringer.Option) *br {
	if len(opts) == 0 {
		return b
	}

	b_ := *b
	b_.apply(opts)
	return &b_
}

func HttpBringer(opts ...bringer.Option) bringer.Bringer {
	b := &br{}
	b.apply(opts)
	return b
}

func (b *br) Bring(ctx context.Context, t thing.Thing, opts ...bringer.Option) (io.ReadCloser, error) {
	l := log.From(ctx).With(slog.String("bringer", "http"))
	b = b.with(opts)

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
	res, err := b.client.Get(t.Url.String())
	if err != nil {
		e := &url.Error{}
		if errors.As(err, &e) {
			return nil, err
		}
		return nil, fmt.Errorf("request GET: %w", err)
	}
	l.Info("response", slog.Int("status", res.StatusCode))
	if l.Enabled(ctx, slog.LevelDebug) {
		l.Debug("response", slog.String("header", fmt.Sprintf("%v", res.Header)))
	}
	if res.StatusCode != http.StatusOK {
		if l.Enabled(ctx, slog.LevelDebug) {
			if data, err := io.ReadAll(res.Body); err == nil {
				l.Debug("response", slog.String("body", string(data)))
			}
		}
		return nil, fmt.Errorf("status not 200: %d", res.StatusCode)
	}

	return res.Body, nil
}

func init() {
	bringer.Register("http", HttpBringer)
	bringer.Register("https", HttpBringer)
}
