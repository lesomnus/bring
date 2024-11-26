package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/lesomnus/bring/config"
	"github.com/lesomnus/bring/entry"
	"github.com/lesomnus/bring/internal/hook"
	"github.com/lesomnus/bring/internal/hooks"
	"github.com/lesomnus/bring/internal/task"
	"github.com/lesomnus/bring/log"
	"github.com/lesomnus/bring/thing"
	"github.com/urfave/cli/v3"
)

func NewCmdBring() *cli.Command {
	dest := ""
	dest_given := false
	dry_run := false

	return &cli.Command{
		Name:   "",
		Hidden: true,

		Flags: []cli.Flag{
			&cli.StringFlag{
				Category: "Bring",

				Name:    "into",
				Aliases: []string{"o"},
				Usage:   "Destination to place things",
				Action: func(ctx context.Context, cmd *cli.Command, v string) error {
					dest_given = true
					return nil
				},

				Destination: &dest,
			},
			&cli.BoolFlag{
				Category: "Bring",

				Name:  "dry-run",
				Usage: "Stops before write things to file",
				Value: dry_run,

				Destination: &dry_run,
			},
		},

		Action: func(ctx context.Context, cmd *cli.Command) error {
			inventory_path := "things.yaml"
			switch cmd.NArg() {
			case 0:
				break
			case 1:
				inventory_path = cmd.Args().Get(0)

			default:
				return fmt.Errorf("expected 0 or 1 argument")
			}
			if !dest_given {
				return errors.New("destination must be given; use --into")
			}
			if dest == "" {
				return errors.New(`destination cannot be empty; give "./" to bring into current directory`)
			}

			c := config.From(ctx)
			l := log.From(ctx)
			l.Info("things will be loaded", slog.String("from", inventory_path))
			l.Info("things will be placed", slog.String("to", dest))

			if info, err := os.Stat(dest); err != nil && !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("open destination: %w", err)
			} else if err == nil && !info.IsDir() {
				return errors.New("destination must be a directory")
			}

			i, err := entry.FromPath(inventory_path)
			if err != nil {
				return fmt.Errorf("read inventory: %w", err)
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

			num_tasks := i.Things.Len()
			num_tasks_done := 0
			num_errors := 0
			err = i.Things.Walk(dest, func(p string, t *thing.Thing) error {
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
