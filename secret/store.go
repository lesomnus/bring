package secret

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"net/url"
	"os"

	"github.com/lesomnus/bring/log"
)

type Store interface {
	Read(ctx context.Context, u url.URL) ([]byte, error)
}

func FromUrl(ctx context.Context, u *url.URL) (Store, error) {
	l := log.From(ctx)

	switch u.Scheme {
	case "":
		fallthrough
	case "file":
		p := u.Path
		l.Debug("use FsStore", slog.String("path", p))
		f, ok := os.DirFS(p).(fs.ReadDirFS)
		if !ok {
			panic("`os.DirFS` must implement `fs.ReadDirFS`")
		}
		return FsStore(f), nil

	default:
		return nil, fmt.Errorf("scheme %s not supported", u.Scheme)
	}
}
