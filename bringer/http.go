package bringer

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/lesomnus/bring/thing"
)

type httpBringer struct{}

func HttpBringer(opts ...Option) Bringer {
	return &httpBringer{}
}

func (b *httpBringer) Bring(ctx context.Context, t thing.Thing) (io.ReadCloser, error) {
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

	res, err := http.Get(t.Url.String())
	if err != nil {
		return nil, fmt.Errorf("request get: %w", err)
	}

	return res.Body, nil
}
