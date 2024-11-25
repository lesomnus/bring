package config

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/lesomnus/bring/bringer"
	"github.com/lesomnus/bring/log"
	"github.com/lesomnus/bring/secret"
)

type SecretConfig struct {
	Enabled bool
	Url     string
}

func (c *SecretConfig) Open(ctx context.Context) (secret.Store, error) {
	l := log.From(ctx)

	if !c.Enabled {
		l.Info("secret store is disabled")
		return secret.NopStore(), nil
	}

	u, err := url.Parse(c.Url)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	return secret.FromUrl(ctx, u)
}

func (c *SecretConfig) OpenTo(ctx context.Context, u url.URL, store *secret.Store) error {
	if *store != nil {
		return nil
	}
	if u.User.Username() == "" {
		return nil
	}

	s, err := c.Open(ctx)
	if err != nil {
		return fmt.Errorf("open secret store: %w", err)
	}

	*store = s
	return nil
}

func (c *SecretConfig) AsOpts(ctx context.Context, store secret.Store, u url.URL) ([]bringer.Option, error) {
	if u.User.Username() == "" {
		return nil, nil
	}
	if store == nil {
		store = secret.NopStore()
	}

	l := log.From(ctx)
	if _, ok := u.User.Password(); ok {
		l.Info("use password", slog.String("source", "URL"))
		return nil, nil
	}

	l.Info("use password", slog.String("source", "store"))
	pw, err := store.Read(ctx, u)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("read secret: %w", err)
	}

	return []bringer.Option{bringer.WithPassword(string(pw))}, nil
}
