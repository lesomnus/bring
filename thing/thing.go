package thing

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/opencontainers/go-digest"
	"gopkg.in/yaml.v3"
)

type Thing struct {
	Url    url.URL
	Digest *digest.Digest
}

func (t *Thing) UnmarshalYAML(n *yaml.Node) error {
	obj := map[string]yaml.Node{}
	if err := n.Decode(&obj); err != nil {
		return err
	}

	if n, ok := obj["url"]; !ok {
		return errors.New("url must be set")
	} else {
		if v, err := parseUrl(&n); err != nil {
			return fmt.Errorf("invalid URL: %w", err)
		} else {
			t.Url = *v
		}
	}
	if n, ok := obj["digest"]; ok {
		if v, err := parseDigest(&n); err != nil {
			return fmt.Errorf("invalid digest: %w", err)
		} else {
			t.Digest = &v
		}
	}

	return nil
}

func parseUrl(n *yaml.Node) (*url.URL, error) {
	if n.Kind != yaml.ScalarNode {
		return nil, errors.New("expected `url` to be a string")
	}
	if strings.TrimSpace(n.Value) == "" {
		return nil, errors.New("url cannot be empty")
	}

	return url.Parse(n.Value)
}

func parseDigest(n *yaml.Node) (digest.Digest, error) {
	if n.Kind != yaml.ScalarNode {
		return "", errors.New("expected `url` to be a string")
	}

	v, err := digest.Parse(n.Value)
	if err != nil {
		return "", fmt.Errorf("%s: %w", n.Value, err)
	}

	return v, nil
}
