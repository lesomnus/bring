package thing

import (
	"context"
	"fmt"
	"io"
	"os"
)

type fileBringer struct {
	t Thing
}

func FileBringer(t Thing) Bringer {
	return &fileBringer{t}
}

func (b *fileBringer) Thing() Thing {
	return b.t
}

func (b *fileBringer) Bring(ctx context.Context) (io.ReadCloser, error) {
	f, err := os.Open(b.t.Url.Path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	return f, nil
}
