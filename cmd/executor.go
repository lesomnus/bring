package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/lesomnus/bring/thing"
	"github.com/opencontainers/go-digest"
)

type ExecuteContext struct {
	context.Context

	N int
	I int

	Path  string
	Thing *thing.Thing
}

type Executor struct {
	Context ExecuteContext
	DryRun  bool
	NewHook func(ctx ExecuteContext) ExecuteHook
}

func (e *Executor) Execute(p string, t *thing.Thing) {
	e.Context.I += 1
	ctx := e.Context
	ctx.Path = p
	ctx.Thing = t

	hook := e.NewHook(ctx)
	hook.OnStart()

	errs := []error{}
	if err := thing.ErrFromUrl(t.Url); err != nil {
		errs = append(errs, err)
	}
	if err := thing.ErrFromDigest(t.Digest); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		hook.OnError(errors.Join(errs...))
		return
	}

	d := filepath.Dir(p)
	if err := os.MkdirAll(d, os.ModePerm); err != nil {
		hook.OnError(fmt.Errorf("mkdir: %w", err))
		return
	}

	if ok, err := e.validate(p, t); err != nil {
		hook.OnError(err)
		return
	} else if ok {
		// Digest matches, skips bringing.
		hook.OnSkip()
		return
	}

	if e.DryRun {
		hook.OnFinish()
		return
	}

	r, err := t.Bring(e.Context)
	if err != nil {
		hook.OnError(fmt.Errorf("bring: %w", err))
		return
	}
	defer r.Close()

	w, err := os.Create(p)
	if err != nil {
		hook.OnError(fmt.Errorf("create: %w", err))
		return
	}
	defer w.Close()

	if _, err := io.Copy(w, r); err != nil {
		hook.OnError(fmt.Errorf("bringing: %w", err))
	}

	hook.OnDone()
	hook.OnFinish()
}

func (e *Executor) validate(p string, t *thing.Thing) (bool, error) {
	algo := t.Digest.Algorithm()
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

	d := digest.NewDigest(algo, hash)
	if d == t.Digest {
		// Digest matches, skips bringing.
		return true, nil
	}

	return false, nil
}
