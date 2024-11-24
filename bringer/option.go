package bringer

import "time"

type Option interface {
	_option()
}

type option struct{}

func (*option) _option() {}

type pwOpt struct {
	option
	v string
}

func WithPassword(v string) Option {
	return &pwOpt{v: v}
}

type dialTimeoutOpt struct {
	option
	v time.Duration
}

func WithDialTimeout(v time.Duration) Option {
	return &dialTimeoutOpt{v: v}
}
