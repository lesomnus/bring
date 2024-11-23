package thing

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/opencontainers/go-digest"
	"gopkg.in/yaml.v3"
)

var bringers = map[string](func(t Thing) Bringer){
	"file":  FileBringer,
	"http":  HttpBringer,
	"https": HttpBringer,
}

type Thing struct {
	Name   string
	Url    *url.URL
	Digest digest.Digest
}

func (t *Thing) Bringer() (Bringer, error) {
	b, ok := bringers[t.Url.Scheme]
	if !ok {
		return nil, fmt.Errorf("schema is not supported: %s", t.Url.Scheme)
	}

	return b(*t), nil
}

func (t *Thing) Bring(ctx context.Context) (io.ReadCloser, error) {
	b, err := t.Bringer()
	if err != nil {
		return nil, err
	}

	return SafeBringer(b).Bring(ctx)
}

func (t *Thing) UnmarshalYAML(n *yaml.Node) error {
	obj := map[string]yaml.Node{}
	if err := n.Decode(&obj); err != nil {
		return err
	}

	if n, ok := obj["name"]; ok {
		t.Name = n.Value
	}
	if n, ok := obj["url"]; ok {
		t.Url = parseUrl(&n)
	}
	if n, ok := obj["digest"]; ok {
		t.Digest = parseDigest(&n)
	}

	return nil
}

func parseUrl(n *yaml.Node) *url.URL {
	v := &url.URL{}
	var err error
	if n.Kind != yaml.ScalarNode {
		return NewFailedUrl(FailInvalid, "expected `url` to be a string")
	}

	v, err = url.Parse(n.Value)
	if err != nil {
		return NewFailedUrl(FailInvalid, err.Error())
	}
	if _, ok := bringers[v.Scheme]; !ok {
		return NewFailedUrl(FailNotSupported, v.Scheme)
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

		d := digest.NewDigestFromHex(o.Algo, o.Value)
		if !d.Algorithm().Available() {
			return NewFailedDigest(FailNotSupported, fmt.Sprintf("algorithm %s not supported", d.Algorithm().String()))
		}
		if err := d.Validate(); err != nil {
			return NewFailedDigest(FailInvalid, err.Error())
		}
		return digest.NewDigestFromHex(o.Algo, o.Value)

	default:
		return NewFailedDigest(FailInvalid, "expected `digest` field to be a string or a map")
	}
}
