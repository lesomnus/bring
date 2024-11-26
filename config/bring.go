package config

import (
	"context"
	"time"

	"github.com/lesomnus/bring/bringer"
)

type BringConfig struct {
	TimeoutBring time.Duration `yaml:"timeout_bring"`
	TimeoutDial  time.Duration `yaml:"timeout_dial"`
}

func (c *BringConfig) ApplyBringTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if c.TimeoutBring == 0 {
		return ctx, func() {}
	}

	return context.WithTimeout(ctx, c.TimeoutBring)
}

func (c *BringConfig) AsOpts() []bringer.Option {
	opts := []bringer.Option{}
	if c.TimeoutDial != 0 {
		opts = append(opts, bringer.WithDialTimeout(c.TimeoutDial))
	}

	return opts
}
