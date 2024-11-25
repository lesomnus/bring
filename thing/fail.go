package thing

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/opencontainers/go-digest"
)

type Fail string

const (
	failPrefix = "!"

	FailInvalid Fail = "!invalid"
)

func NewFailedUrl(f Fail, msg string) *url.URL {
	return &url.URL{Scheme: string(f), Fragment: msg}
}

func FailFromUrl(u *url.URL) (Fail, bool) {
	if !strings.HasPrefix(u.Scheme, failPrefix) {
		return "", false
	}

	return Fail(u.Scheme), true
}

func ErrFromUrl(u *url.URL) error {
	if !strings.HasPrefix(u.Scheme, failPrefix) {
		return nil
	}

	return fmt.Errorf("url: %s", u.Fragment)
}

func NewFailedDigest(f Fail, msg string) digest.Digest {
	return digest.NewDigestFromHex(string(f), msg)
}

func FailFromDigest(d digest.Digest) (Fail, bool) {
	if !strings.HasPrefix(d.Algorithm().String(), failPrefix) {
		return "", false
	}

	return Fail(d.Algorithm()), true
}

func ErrFromDigest(d digest.Digest) error {
	if !strings.HasPrefix(string(d), failPrefix) {
		return nil
	}

	return fmt.Errorf("digest: %s", d.Encoded())
}
