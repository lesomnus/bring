package option

import "time"

type Option interface {
	_option()
}

type OptionTag struct{}

func (*OptionTag) _option() {}

type PwOption struct {
	OptionTag
	Value string
}

type DialTimeoutOpt struct {
	OptionTag
	Value time.Duration
}
