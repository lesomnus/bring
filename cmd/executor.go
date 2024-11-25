package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/lesomnus/bring/bringer"
	"github.com/lesomnus/bring/internal/hook"
	"github.com/lesomnus/bring/internal/task"
	"github.com/lesomnus/bring/secret"
	"github.com/opencontainers/go-digest"
)

type executor struct {
	Secret secret.Store

	DryRun  bool
	NewHook func(ctx context.Context, t task.Task) hook.Hook
}

func (e *executor) Execute(ctx context.Context, t task.Task) {
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

	opts := []bringer.Option{}
	opts = append(opts, t.BringConfig.AsOpts()...)
	if pw, err := e.Secret.Read(ctx, t.Thing.Url); err != nil && !errors.Is(err, os.ErrNotExist) {
		hook.OnError(fmt.Errorf("read secret: %w", err))
		return
	} else {
		opts = append(opts, bringer.WithPassword(string(pw)))
	}
	if t.BringConfig.DialTimeout != 0 {
		opts = append(opts, bringer.WithDialTimeout(t.BringConfig.DialTimeout))
	}

	ctx, cancel := t.BringConfig.ApplyTimeout(ctx)
	defer cancel()

	r, err := b.Bring(ctx, t.Thing, opts...)
	if err != nil {
		hook.OnError(fmt.Errorf("bring: %w", err))
		return
	} else {
		defer r.Close()
	}

	if e.DryRun {
		hook.OnDone(nil)
		return
	}

	hook.OnDone(r)
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
