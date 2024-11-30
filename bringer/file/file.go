package file

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/lesomnus/bring/bringer"
	"github.com/lesomnus/bring/log"
	"github.com/lesomnus/bring/thing"
)

type br struct{}

func FileBringer(opts ...bringer.Option) bringer.Bringer {
	return &br{}
}

func (b *br) Bring(ctx context.Context, t thing.Thing, opts ...bringer.Option) (io.ReadCloser, error) {
	l := log.From(ctx).With(slog.String("bringer", "file"))

	p := fmt.Sprintf("%s%s", t.Url.Host, t.Url.Path)
	l.Info("open", slog.String("path", p))

	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	return f, nil
}

func init() {
	bringer.Register("", FileBringer)
	bringer.Register("file", FileBringer)
}
