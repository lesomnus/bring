package thing

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type httpBringer struct {
	t Thing
}

func HttpBringer(t Thing) Bringer {
	return &httpBringer{t}
}

func (b *httpBringer) Thing() Thing {
	return b.t
}

func (b *httpBringer) Bring(ctx context.Context) (io.ReadCloser, error) {
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

	res, err := http.Get(b.t.Url.String())
	if err != nil {
		return nil, fmt.Errorf("request get: %w", err)
	}

	return res.Body, nil
}
