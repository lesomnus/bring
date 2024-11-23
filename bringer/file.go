package bringer

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/lesomnus/bring/thing"
)

type fileBringer struct{}

func FileBringer(opts ...Option) Bringer {
	return &fileBringer{}
}

func (b *fileBringer) Bring(ctx context.Context, t thing.Thing) (io.ReadCloser, error) {
	f, err := os.Open(t.Url.Path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	return f, nil
}
