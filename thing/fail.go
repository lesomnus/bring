package thing

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/opencontainers/go-digest"
)

const (
	FailPrefix       = "!"
	FailInvalid      = "!invalid"
	FailNotSupported = "!not supported"
)

func ErrFromUrl(u *url.URL) error {
	if !strings.HasPrefix(u.Scheme, FailPrefix) {
		return nil
	}

	return fmt.Errorf("url: %s", u.Fragment)
}

func ErrFromDigest(d digest.Digest) error {
	if !strings.HasPrefix(d.Algorithm().String(), FailPrefix) {
		return nil
	}

	return fmt.Errorf("digest: %s", d.Encoded())
}
