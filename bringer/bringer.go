package bringer

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/lesomnus/bring/thing"
)

type Bringer interface {
	Bring(ctx context.Context, t thing.Thing, opts ...Option) (io.ReadCloser, error)
}

func FromUrl(u url.URL, opts ...Option) (Bringer, error) {
	bf, ok := bringers[u.Scheme]
	if !ok {
		return nil, fmt.Errorf("scheme not supported: %s", u.Scheme)
	}

	return bf(opts...), nil
}
