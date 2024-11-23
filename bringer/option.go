package bringer

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

// TODO: TCP timeout
