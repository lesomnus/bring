package http

import (
	"net/http"

	"github.com/lesomnus/bring/bringer"
	"github.com/lesomnus/bring/bringer/option"
)

type transportOption struct {
	option.OptionTag
	Value http.RoundTripper
}

func WithTransport(v http.RoundTripper) bringer.Option {
	return &transportOption{Value: v}
}
