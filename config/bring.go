package config

import (
	"context"
	"time"

	"github.com/lesomnus/bring/bringer"
)

type BringConfig struct {
	BringTimeout time.Duration `yaml:"bring_timeout"`
	DialTimeout  time.Duration `yaml:"dial_timeout"`
}

func (c *BringConfig) ApplyBringTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if c.BringTimeout == 0 {
		return ctx, func() {}
	}

	return context.WithTimeout(ctx, c.BringTimeout)
}

func (c *BringConfig) AsOpts() []bringer.Option {
	opts := []bringer.Option{}
	if c.DialTimeout != 0 {
		opts = append(opts, bringer.WithDialTimeout(c.DialTimeout))
	}

	return opts
}
