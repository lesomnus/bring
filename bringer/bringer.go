package bringer

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/url"

	"github.com/lesomnus/bring/thing"
)

type Bringer interface {
	Bring(ctx context.Context, t thing.Thing, opts ...Option) (io.ReadCloser, error)
}

func FromUrl(u url.URL, opts ...Option) (Bringer, error) {
	bf, ok := bringers[u.Scheme]
	if !ok {
		return nil, fmt.Errorf("scheme %s not supported", u.Scheme)
	}

	return bf(opts...), nil
}

func name(name string) slog.Attr {
	return slog.String("bringer", name)
}
