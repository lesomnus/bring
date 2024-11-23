package bringer

import (
	"context"
	"fmt"
	"io"

	"github.com/lesomnus/bring/thing"
)

type smbBringer struct {
	password string
}

func SmbBringer(opts ...Option) Bringer {
	b := &smbBringer{}
	for _, opt := range opts {
		switch o := opt.(type) {
		case (*pwOpt):
			b.password = o.v
		}
	}

	return b
}

func (b *smbBringer) Bring(ctx context.Context, t thing.Thing) (io.ReadCloser, error) {
	return nil, fmt.Errorf("not implemented")
}
