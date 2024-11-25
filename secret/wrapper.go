package secret

import (
	"context"
	"log/slog"
	"net/url"
	"os"

	"github.com/lesomnus/bring/log"
)

type wrapper struct {
	s Store
}

func (w wrapper) Read(ctx context.Context, u url.URL) ([]byte, error) {
	if u.User.Username() == "" {
		return nil, os.ErrNotExist
	}

	l := log.From(ctx)
	if v, ok := u.User.Password(); ok {
		l.Info("read password", slog.String("source", "URL"))
		return []byte(v), nil
	}

	return w.s.Read(ctx, u)
}
