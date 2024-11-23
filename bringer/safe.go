package bringer

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/lesomnus/bring/log"
	"github.com/lesomnus/bring/thing"
	"github.com/opencontainers/go-digest"
)

type safeBringer struct {
	Bringer
}

func SafeBringer(b Bringer) Bringer {
	return &safeBringer{b}
}

func (b *safeBringer) Bring(ctx context.Context, t thing.Thing) (io.ReadCloser, error) {
	l := log.From(ctx).With(name("safe"))

	algorithm := t.Digest.Algorithm()
	if !algorithm.Available() {
		return nil, fmt.Errorf("unknown type of digest: %s", string(algorithm))
	}

	r, err := b.Bringer.Bring(ctx, t)
	if err != nil {
		return nil, err
	}

	if s, ok := r.(io.ReadSeekCloser); ok {
		// We can seek to start after calculating a digest rather than copy.
		l.Debug("source is seek-able; seek start")
		return b.finalizeSeeker(t.Digest, s)
	}

	sink := &bytes.Buffer{} // TODO: it can be a file.
	if err := b.copyWithVerify(t.Digest, r, sink); err != nil {
		return nil, err
	}

	return io.NopCloser(sink), nil
}

func (b *safeBringer) finalizeSeeker(d digest.Digest, r io.ReadSeekCloser) (io.ReadCloser, error) {
	if err := b.copyWithVerify(d, r, io.Discard); err != nil {
		return nil, err
	}
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("seek to start: %w", err)
	}

	return r, nil
}

func (b *safeBringer) copyWithVerify(d digest.Digest, r io.Reader, sink io.Writer) error {
	algo := d.Algorithm()
	hash := algo.Hash()
	w := io.MultiWriter(hash, sink)
	if _, err := io.Copy(w, r); err != nil {
		// Note that `hash` never returns an error.
		return fmt.Errorf("copy: %w", err)
	}

	digest := digest.NewDigest(algo, hash)
	if digest != d {
		return fmt.Errorf("digest mismatch: %s != %s", digest, d)
	}

	return nil
}
