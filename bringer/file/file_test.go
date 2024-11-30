package file_test

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/lesomnus/bring/bringer/file"
	"github.com/lesomnus/bring/thing"
	"github.com/stretchr/testify/require"
)

func TestFileBringer(t *testing.T) {
	b := file.FileBringer()

	d := t.TempDir()
	p := filepath.Join(d, "foo")
	data := []byte("bar")
	err := os.WriteFile(p, data, 0o644)
	require.NoError(t, err)

	t.Run("absolute path without schema", func(t *testing.T) {
		require := require.New(t)

		u, err := url.Parse(p)
		if err != nil {
			require.NoError(err)
		}

		f, err := b.Bring(context.Background(), thing.Thing{Url: *u})
		if err == nil {
			defer f.Close()
		}
		require.NoError(err)

		v, err := io.ReadAll(f)
		require.NoError(err)
		require.Equal(v, data)
	})
	t.Run("absolute path with schema", func(t *testing.T) {
		require := require.New(t)

		u, err := url.Parse(fmt.Sprintf("file://%s", p))
		if err != nil {
			require.NoError(err)
		}

		f, err := b.Bring(context.Background(), thing.Thing{Url: *u})
		if err == nil {
			defer f.Close()
		}
		require.NoError(err)

		v, err := io.ReadAll(f)
		require.NoError(err)
		require.Equal(v, data)
	})
	t.Run("relative path", func(t *testing.T) {
		require := require.New(t)

		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		if err := os.Chdir(d); err != nil {
			panic(err)
		}
		defer os.Chdir(wd)

		u, err := url.Parse("./foo")
		if err != nil {
			require.NoError(err)
		}

		f, err := b.Bring(context.Background(), thing.Thing{Url: *u})
		if err == nil {
			defer f.Close()
		}
		require.NoError(err)

		v, err := io.ReadAll(f)
		require.NoError(err)
		require.Equal(v, data)
	})
}
