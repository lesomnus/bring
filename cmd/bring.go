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
	"github.com/lesomnus/bring/thing"
	"github.com/urfave/cli/v3"
)

func NewCmdBring() *cli.Command {
	executor := &executor{
		DryRun: false,
	}

	return &cli.Command{
		Name:  "bring",
		Usage: "Bring things.",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "Stops before write things to file",
				Value: executor.DryRun,

				Destination: &executor.DryRun,
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

			conf := config.From(ctx)
			if dest != "" {
				conf.Dest = dest
			}
			if conf.Dest == "" {
				return fmt.Errorf("destination must be specified in the config file or given by argument")
			}

			var err error
			executor.Secret, err = conf.Secret.Open(ctx)
			if err != nil {
				return fmt.Errorf("open secret store: %w", err)
			}

			num_errors := 0
			executor.NewHook = func(ctx context.Context, t task.Task) hook.Hook {
				return hook.Tie(
					&sinkHookMw{D: t.Dest},
					hook.Forward(
						hook.Join(
							&countErrHook{n: &num_errors},
							&hooks.PrintHook{T: t, O: os.Stdout},
						),
					),
				)
			}

			job := task.Job{
				NumTasks: conf.Things.Len(),
			}
			i := 0
			conf.Things.Walk(conf.Dest, func(p string, t *thing.Thing) {
				task := task.Task{
					Thing: *t,

					BringConfig: conf.Each,

					Job:   job,
					Order: i,
					Dest:  p,
				}

				executor.Execute(ctx, task)
				i++
			})

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

type countErrHook struct {
	hook.NopHook
	n *int
}

func (h *countErrHook) OnError(err error) {
	*h.n++
}
