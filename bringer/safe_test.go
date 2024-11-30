package bringer_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/lesomnus/bring/bringer"
	"github.com/lesomnus/bring/thing"
	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/require"
)

type mockBringer struct {
	r io.ReadCloser
}

func (b *mockBringer) Bring(ctx context.Context, t thing.Thing, opts ...bringer.Option) (io.ReadCloser, error) {
	return b.r, nil
}

func TestSafeBringer(t *testing.T) {
	t.Run("expects a digest", func(t *testing.T) {
		require := require.New(t)

		b := bringer.SafeBringer(&mockBringer{})
		_, err := b.Bring(context.Background(), thing.Thing{
			Digest: nil,
		})
		require.ErrorContains(err, "no digest")
	})
	t.Run("expects a valid digest", func(t *testing.T) {
		require := require.New(t)

		b := bringer.SafeBringer(&mockBringer{})
		d := digest.Digest("invalid:foo")
		_, err := b.Bring(context.Background(), thing.Thing{
			Digest: &d,
		})
		require.ErrorContains(err, "invalid")
	})
	t.Run("expects digest be matched", func(t *testing.T) {
		require := require.New(t)

		data := []byte("Royale with Cheese")
		f := &bytes.Buffer{}
		f.Write(data)

		h := digest.SHA256.Hash()
		h.Write([]byte("Le Big Mac"))
		d := digest.NewDigest(digest.SHA256, h)

		b := bringer.SafeBringer(&mockBringer{r: io.NopCloser(f)})
		_, err := b.Bring(context.Background(), thing.Thing{
			Digest: &d,
		})
		require.ErrorContains(err, "mismatch")
	})
	t.Run("returns original reader if the reader is seek-able", func(t *testing.T) {
		require := require.New(t)

		p := filepath.Join(t.TempDir(), "burger")
		data := []byte("Royale with Cheese")
		err := os.WriteFile(p, data, 0o644)
		require.NoError(err)

		f, err := os.Open(p)
		if err == nil {
			defer f.Close()
		}
		require.NoError(err)

		h := digest.SHA256.Hash()
		h.Write(data)
		d := digest.NewDigest(digest.SHA256, h)

		b := bringer.SafeBringer(&mockBringer{r: f})
		v, err := b.Bring(context.Background(), thing.Thing{
			Digest: &d,
		})
		require.NoError(err)
		require.Same(v, f)
	})
}
