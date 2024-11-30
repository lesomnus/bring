package cmd_test

import (
	"bytes"
	"context"
	_ "crypto/sha256"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/lesomnus/bring/bringer/file"
	. "github.com/lesomnus/bring/cmd"
	"github.com/lesomnus/bring/config"
	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/require"
)

func TestCmdDigest(t *testing.T) {
	t.Run("schema not supported", func(t *testing.T) {
		require := require.New(t)

		c := config.New()
		ctx := config.Into(context.Background(), c)

		cmd := NewCmdDigest()
		err := cmd.Run(ctx, []string{"", "foo://bar"})
		require.ErrorContains(err, "not supported")
	})
	t.Run("prints digest of the resource", func(t *testing.T) {
		require := require.New(t)

		data := []byte("Royale with Cheese")
		p := filepath.Join(t.TempDir(), "burger")
		err := os.WriteFile(p, data, 0o644)
		require.NoError(err)

		h := digest.SHA256.Hash()
		h.Write(data)
		d := digest.NewDigest(digest.SHA256, h)

		c := config.New()
		ctx := config.Into(context.Background(), c)

		b := &bytes.Buffer{}
		cmd := NewCmdDigest()
		cmd.Writer = b
		err = cmd.Run(ctx, []string{"", p})
		require.NoError(err)

		s := b.String()
		s = strings.TrimSpace(s)
		t.Log(s)

		v, err := digest.Parse(s)
		require.NoError(err)
		require.Equal(v, d)
	})
}
