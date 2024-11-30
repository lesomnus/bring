package secret

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/lesomnus/bring/log"
)

// TODO: cache?
type fsStore struct {
	fs fs.ReadDirFS
}

func FsStore(fs fs.ReadDirFS) Store {
	return &wrapper{s: &fsStore{fs: fs}}
}

func (s fsStore) Read(ctx context.Context, u url.URL) ([]byte, error) {
	// Follow the given path down to the deepest matching directory.
	// If there are no more matching subdirectories, look for a file that matches the username in the URL.
	// If the file isn't found in the current directory, move up to the parent directories until the file is located.

	p := filepath.Join(u.Scheme, u.Hostname())
	ds, err := fs.ReadDir(s.fs, p)
	if err != nil {
		return nil, fmt.Errorf(`read dir "%s": %w`, p, err)
	}

	// Holds directory entires of the visited directories.
	dss := [][]fs.DirEntry{ds}

	var d fs.DirEntry

	entries := strings.Split(filepath.Clean(u.Path), "/")
	if len(entries) > 0 && entries[0] == "" {
		entries = entries[1:]
	}
	for _, e := range entries {
		i := slices.IndexFunc(ds, func(d_ fs.DirEntry) bool {
			return d_.Name() == e
		})
		if i < 0 {
			break
		}

		d = ds[i]
		if !d.IsDir() {
			break
		}

		p_next := filepath.Join(p, e)
		ds_next, err := s.fs.ReadDir(p_next)
		if err != nil {
			break
		}

		p = p_next
		ds = ds_next
		dss = append(dss, ds)
	}

	name := "@" + u.User.Username()
	for _, ds := range slices.Backward(dss) {
		i := slices.IndexFunc(ds, func(d_ fs.DirEntry) bool {
			return d_.Name() == name
		})
		if i > -1 {
			p = filepath.Join(p, name)
			break
		}

		p = filepath.Dir(p)
	}
	if p == u.Scheme || p == "." {
		return nil, os.ErrNotExist
	}

	l := log.From(ctx)
	l.Info("read password", slog.String("path", p))
	return fs.ReadFile(s.fs, p)
}
