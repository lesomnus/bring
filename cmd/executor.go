package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/lesomnus/bring/bringer"
	"github.com/lesomnus/bring/internal/hook"
	"github.com/lesomnus/bring/internal/task"
	"github.com/lesomnus/bring/log"
	"github.com/lesomnus/bring/secret"
	"github.com/opencontainers/go-digest"
)

type executor struct {
	Secret secret.Store

	DryRun  bool
	NewHook func(ctx context.Context, t task.Task) hook.Hook
}

func (e *executor) Execute(ctx context.Context, t task.Task) {
	// TODO: move to hook

	l := log.From(ctx)
	l.Info("start",
		slog.String("from", t.Thing.Url.Redacted()),
		slog.String("to", t.Dest),
	)

	hook := e.NewHook(ctx, t)
	hook.OnStart()
	defer hook.OnFinish()

	if err := t.Thing.Validate(); err != nil {
		hook.OnError(fmt.Errorf("invalid thing: %w", err))
		return
	}

	b, err := bringer.FromUrl(t.Thing.Url)
	if err != nil {
		hook.OnError(fmt.Errorf("get bringer: %w", err))
		return
	} else {
		b = bringer.SafeBringer(b)
	}

	if ok, err := e.validate(t.Dest, t.Thing.Digest); err != nil {
		hook.OnError(err)
		return
	} else if ok {
		// Digest matches, skips bringing.
		hook.OnSkip()
		return
	}

	ctx_bring := ctx
	opts := []bringer.Option{}
	if t.Thing.Url.User.Username() != "" {
		if _, ok := t.Thing.Url.User.Password(); ok {
			l.Info("use password", slog.String("source", "URL"))
		} else {
			l.Info("use password", slog.String("source", "store"))

			// TODO: need timeout?
			pw, err := e.Secret.Read(ctx, t.Thing.Url)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				hook.OnError(fmt.Errorf("read secret: %w", err))
				return
			}

			opts = append(opts, bringer.WithPassword(string(pw)))
		}
	}
	if t.BringConfig.BringTimeout != 0 {
		c, cancel := context.WithTimeout(ctx_bring, t.BringConfig.BringTimeout)
		ctx_bring = c
		defer cancel()
	}
	if t.BringConfig.DialTimeout != 0 {
		opts = append(opts, bringer.WithDialTimeout(t.BringConfig.DialTimeout))
	}
	r, err := b.Bring(ctx_bring, t.Thing, opts...)
	if err != nil {
		hook.OnError(fmt.Errorf("bring: %w", err))
		return
	}
	defer r.Close()

	if e.DryRun {
		return
	}

	d := filepath.Dir(t.Dest)
	if err := os.MkdirAll(d, os.ModePerm); err != nil {
		hook.OnError(fmt.Errorf("mkdir: %w", err))
		return
	}

	w, err := os.Create(t.Dest)
	if err != nil {
		hook.OnError(fmt.Errorf("create: %w", err))
		return
	}
	defer w.Close()

	if _, err := io.Copy(w, r); err != nil {
		hook.OnError(fmt.Errorf("bringing: %w", err))
		return
	}

	hook.OnDone()
}

func (e *executor) validate(p string, d digest.Digest) (bool, error) {
	algo := d.Algorithm()
	hash := algo.Hash()

	f, err := os.Open(p)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return false, fmt.Errorf("open: %w", err)
		}

		return false, nil
	}
	defer f.Close()

	if _, err := io.Copy(hash, f); err != nil {
		return false, fmt.Errorf("copy to calculate hash sum: %w", err)
	}

	d_ := digest.NewDigest(algo, hash)
	if d_ == d {
		// Digest matches, skips bringing.
		return true, nil
	}

	return false, nil
}
