package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/lesomnus/bring/config"
	"github.com/lesomnus/bring/internal/hook"
	"github.com/lesomnus/bring/internal/hooks"
	"github.com/lesomnus/bring/internal/task"
	"github.com/lesomnus/bring/log"
	"github.com/lesomnus/bring/thing"
	"github.com/urfave/cli/v3"
)

func NewCmdBring() *cli.Command {
	dry_run := false

	return &cli.Command{
		Name:  "bring",
		Usage: "Bring things.",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "Stops before write things to file",
				Value: dry_run,

				Destination: &dry_run,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			dest := ""
			switch cmd.NArg() {
			case 0:
				return fmt.Errorf("path to config file must be given")
			case 1:
				break
			case 2:
				dest = cmd.Args().Get(1)

			default:
				return fmt.Errorf("expected 1 or 2 arguments")
			}

			c := config.From(ctx)
			if dest != "" {
				c.Dest = dest
			}
			if c.Dest == "" {
				return fmt.Errorf("destination must be specified in the config file or given by argument")
			}

			exe := executor{
				NewHook: func(ctx context.Context, t task.Task) hook.Hook {
					hs := []hook.Mw{}
					if !dry_run {
						hs = append(hs, &sinkHookMw{D: t.Dest})
					}
					hs = append(hs, hook.Forward(hook.Join(
						&hooks.LogHook{T: t, L: log.From(ctx)},
						&hooks.PrintHook{T: t, O: os.Stdout},
					)))

					return hook.Tie(hs...)
				},
			}

			num_tasks := c.Things.Len()
			num_tasks_done := 0
			num_errors := 0
			err := c.Things.Walk(c.Dest, func(p string, t *thing.Thing) error {
				task := task.Task{
					Thing: *t,

					BringConfig: c.Each,

					Job:   task.Job{NumTasks: num_tasks},
					Order: num_tasks_done,
					Dest:  p,
				}

				if err := c.Secret.OpenTo(ctx, t.Url, &exe.Secret); err != nil {
					return err
				}

				ctx, cancel := c.Each.ApplyBringTimeout(ctx)
				defer cancel()

				if r, err := exe.Run(ctx, task); err != nil {
					num_errors++
				} else if r != nil {
					defer r.Close()
				}
				num_tasks_done++

				return nil
			})
			if err != nil {
				return err
			}
			if num_errors > 0 {
				return cli.Exit("failed to bring some of things", 1)
			}

			return nil
		},
	}
}

type sinkHookMw struct {
	hook.NopMw
	D string
}

func (h *sinkHookMw) sink(r io.Reader) error {
	d := filepath.Dir(h.D)
	if err := os.MkdirAll(d, os.ModePerm); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	w, err := os.Create(h.D)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer w.Close()

	if _, err := io.Copy(w, r); err != nil {
		return fmt.Errorf("bringing: %w", err)
	}

	return nil
}

func (h *sinkHookMw) OnDone(next hook.Hook, r io.Reader) {
	if err := h.sink(r); err != nil {
		next.OnError(err)
	} else {
		next.OnDone(r)
	}
}
