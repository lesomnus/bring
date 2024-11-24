package config

import (
	"context"
	"fmt"
	"net/url"

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
