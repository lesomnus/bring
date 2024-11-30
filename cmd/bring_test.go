package cmd_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	. "github.com/lesomnus/bring/cmd"
	"github.com/lesomnus/bring/config"
	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func TestCmdBring(t *testing.T) {
	c := config.New()
	ctx := config.Into(context.Background(), c)

	t.Run("destination not given", func(t *testing.T) {
		require := require.New(t)

		cmd := NewCmdBring()
		err := cmd.Run(ctx, []string{""})
		require.ErrorContains(err, "destination")
	})
	t.Run("destination not be an empty string", func(t *testing.T) {
		require := require.New(t)

		cmd := NewCmdBring()
		err := cmd.Run(ctx, []string{"", "--into", ""})
		require.ErrorContains(err, "destination")
	})
	t.Run("cannot bring into where not the directory", func(t *testing.T) {
		require := require.New(t)

		cmd := NewCmdBring()
		err := cmd.Run(ctx, []string{"", "--into", "/dev/null"})
		require.ErrorContains(err, "directory")
	})
	t.Run("inventory must be exist", func(t *testing.T) {
		require := require.New(t)

		cmd := NewCmdBring()
		err := cmd.Run(ctx, []string{"", "not-exists", "--into", t.TempDir()})
		require.ErrorContains(err, "inventory")
	})
	t.Run("exit with non-zero code if bring some of things are failed", func(t *testing.T) {
		require := require.New(t)

		p := filepath.Join(t.TempDir(), "things.yaml")
		err := os.WriteFile(p, []byte(`
things:
  something:
    url: /not-exists
    digest: sha256:12794390cce7d0682ffc783c785e4282305684431b30b29ed75c224da24035b4
`), 0o644)
		require.NoError(err)

		cmd := NewCmdBring()
		cmd.ExitErrHandler = func(ctx context.Context, c *cli.Command, err error) {}
		err = cmd.Run(ctx, []string{"", p, "--into", t.TempDir()})
		require.ErrorContains(err, "bring")
	})
	t.Run("bring things into directory", func(t *testing.T) {
		require := require.New(t)

		data := []byte("Royale with Cheese")
		p_file := filepath.Join(t.TempDir(), "burger")
		err := os.WriteFile(p_file, data, 0o644)
		require.NoError(err)

		h := digest.SHA256.Hash()
		h.Write(data)
		d := digest.NewDigest(digest.SHA256, h)

		inventory := fmt.Sprintf(`
things:
  something:
    url: %s
    digest: %s
`,
			p_file,
			d,
		)
		p_inventory := filepath.Join(t.TempDir(), "things.yaml")
		err = os.WriteFile(p_inventory, []byte(inventory), 0o644)
		require.NoError(err)

		cmd := NewCmdBring()
		cmd.ExitErrHandler = func(ctx context.Context, c *cli.Command, err error) {}
		p_dest := t.TempDir()
		err = cmd.Run(ctx, []string{"", p_inventory, "--into", p_dest})
		require.NoError(err)

		v, err := os.ReadFile(filepath.Join(p_dest, "something"))
		require.NoError(err)
		require.Equal(v, data)
	})
}
