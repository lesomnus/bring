package thing

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/opencontainers/go-digest"
	"gopkg.in/yaml.v3"
)

type Thing struct {
	Url    url.URL
	Digest digest.Digest
}

func (t *Thing) UnmarshalYAML(n *yaml.Node) error {
	obj := map[string]yaml.Node{}
	if err := n.Decode(&obj); err != nil {
		return err
	}

	if n, ok := obj["url"]; ok {
		t.Url = *parseUrl(&n)
	}
	if n, ok := obj["digest"]; ok {
		t.Digest = parseDigest(&n)
	}

	return nil
}

func (t *Thing) Validate() error {
	errs := []error{}
	if err := ErrFromUrl(&t.Url); err != nil {
		errs = append(errs, err)
	}
	if err := ErrFromDigest(t.Digest); err != nil {
		errs = append(errs, err)
	}
	if t.Digest != "" {
		if err := t.Digest.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("invalid digest: %w", err))
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func parseUrl(n *yaml.Node) *url.URL {
	if n.Kind != yaml.ScalarNode {
		return NewFailedUrl(FailInvalid, "expected `url` to be a string")
	}

	v, err := url.Parse(n.Value)
	if err != nil {
		return NewFailedUrl(FailInvalid, err.Error())
	}

	return v
}

func parseDigest(n *yaml.Node) digest.Digest {
	switch n.Kind {
	case yaml.ScalarNode:
		return digest.Digest(n.Value)

	case yaml.MappingNode:
		o := struct {
			Algo  string
			Value string
		}{}
		if err := n.Decode(&o); err != nil {
			return NewFailedDigest(FailInvalid, err.Error())
		}

		return digest.NewDigestFromHex(o.Algo, o.Value)

	default:
		return NewFailedDigest(FailInvalid, "expected `digest` field to be a string or a map")
	}
}
