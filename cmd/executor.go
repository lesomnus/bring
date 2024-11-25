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

func (e *executor) secret() secret.Store {
	if e.Secret != nil {
		return e.Secret
	}

	return secret.NopStore()
}

func (e *executor) Execute(ctx context.Context, t task.Task) (io.ReadCloser, error) {
	hook := e.NewHook(ctx, t)
	hook.OnStart()
	defer hook.OnFinish()

	with_validate := !(t.Dest == "" || t.Thing.Digest == "")

	if err := t.Thing.Validate(); err != nil {
		err = fmt.Errorf("invalid thing: %w", err)
		hook.OnError(err)
		return nil, err
	}

	b, err := bringer.FromUrl(t.Thing.Url)
	if err != nil {
		err = fmt.Errorf("get bringer: %w", err)
		hook.OnError(err)
		return nil, err
	} else if with_validate {
		b = bringer.SafeBringer(b)
	}

	if with_validate {
		if ok, err := e.validate(t.Dest, t.Thing.Digest); err != nil {
			hook.OnError(err)
			return nil, err
		} else if ok {
			// Digest matches, skips bringing.
			hook.OnSkip()
			return nil, nil
		}
	}

	opts := []bringer.Option{}
	opts = append(opts, t.BringConfig.AsOpts()...)
	if pw, err := e.secret().Read(ctx, t.Thing.Url); err != nil && !errors.Is(err, os.ErrNotExist) {
		err = fmt.Errorf("read secret: %w", err)
		hook.OnError(err)
		return nil, err
	} else {
		opts = append(opts, bringer.WithPassword(string(pw)))
	}

	r, err := b.Bring(ctx, t.Thing, opts...)
	if err != nil {
		err := fmt.Errorf("bring: %w", err)
		hook.OnError(err)
		return nil, err
	}

	if e.DryRun {
		hook.OnDone(nil)
		return r, nil
	}

	hook.OnDone(r)
	return r, nil
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
