package bringer

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/lesomnus/bring/thing"
)

type Bringer interface {
	Bring(ctx context.Context, t thing.Thing) (io.ReadCloser, error)
}

func FromUrl(u *url.URL) (Bringer, error) {
	bf, ok := bringers[u.Scheme]
	if !ok {
		return nil, fmt.Errorf("scheme %s not supported", u.Scheme)
	}

	return bf(), nil
}
