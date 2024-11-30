package bringer

import (
	"time"

	"github.com/lesomnus/bring/bringer/option"
)

type Option = option.Option

func WithPassword(v string) Option {
	return &option.PwOption{Value: v}
}

func WithDialTimeout(v time.Duration) Option {
	return &option.DialTimeoutOpt{Value: v}
}
