package secret_test

import (
	"context"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/lesomnus/bring/secret"
	"github.com/stretchr/testify/require"
)

func TestFsStore(t *testing.T) {
	test := func(f func(require *require.Assertions, d string, s secret.Store)) func(t *testing.T) {
		return func(t *testing.T) {
			d := t.TempDir()
			fs := os.DirFS(d).(fs.ReadDirFS)
			s := secret.FsStore(fs)
			f(require.New(t), d, s)
		}
	}

	u, err := url.Parse("http://qux@example.com/foo/bar/baz")
	if err != nil {
		panic(err)
	}

	t.Run("read from URL", test(func(require *require.Assertions, d string, s secret.Store) {
		u := *u
		pw := "cheese"
		u.User = url.UserPassword("royale", pw)

		v, err := s.Read(context.Background(), u)
		require.NoError(err)
		require.Equal([]byte(pw), v)
	}))
	t.Run("read exact", test(func(require *require.Assertions, d string, s secret.Store) {
		p := filepath.Join(d, u.Scheme, u.Host, u.Path)
		err := os.MkdirAll(p, os.ModePerm)
		require.NoError(err)

		data := []byte("secret")
		err = os.WriteFile(filepath.Join(p, "@"+u.User.Username()), data, 0o644)
		require.NoError(err)

		v, err := s.Read(context.Background(), *u)
		require.NoError(err)
		require.Equal(v, data)
	}))
	t.Run("read best match", test(func(require *require.Assertions, d string, s secret.Store) {
		p := filepath.Join(d, u.Scheme, u.Host, u.Path)
		err := os.MkdirAll(p, os.ModePerm)
		require.NoError(err)

		data := []byte("secret")
		err = os.WriteFile(filepath.Join(d, u.Scheme, u.Host, "foo", "@"+u.User.Username()), data, 0o644)
		require.NoError(err)

		v, err := s.Read(context.Background(), *u)
		require.NoError(err)
		require.Equal(v, data)
	}))
	t.Run("not exists", test(func(require *require.Assertions, d string, s secret.Store) {
		_, err := s.Read(context.Background(), *u)
		require.ErrorIs(err, os.ErrNotExist)
	}))
}
