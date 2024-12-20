package secret

import (
	"context"
	"net/url"
	"os"
)

type nopStore struct{}

func NopStore() Store {
	return &wrapper{s: &nopStore{}}
}

func (s nopStore) Read(ctx context.Context, u url.URL) ([]byte, error) {
	return nil, os.ErrNotExist
}
