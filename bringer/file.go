package bringer

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/lesomnus/bring/log"
	"github.com/lesomnus/bring/thing"
)

type fileBringer struct{}

func FileBringer(opts ...Option) Bringer {
	return &fileBringer{}
}

func (b *fileBringer) Bring(ctx context.Context, t thing.Thing, opts ...Option) (io.ReadCloser, error) {
	l := log.From(ctx).With(name("file"))

	p := fmt.Sprintf("%s%s", t.Url.Host, t.Url.Path)
	l.Info("open", slog.String("path", p))

	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	return f, nil
}
