package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/lesomnus/bring/config"
	"github.com/lesomnus/bring/internal/hook"
	"github.com/lesomnus/bring/internal/task"
	"github.com/lesomnus/bring/thing"
	"github.com/urfave/cli/v2"
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
		Action: func(c *cli.Context) error {
			dest := ""
			switch c.NArg() {
			case 0:
				return fmt.Errorf("path to config file must be given")
			case 1:
				break
			case 2:
				dest = c.Args().Get(1)

			default:
				return fmt.Errorf("expected 1 or 2 arguments")
			}

			conf_path := c.Args().First()
			conf, err := config.LoadFromFilepath(conf_path)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			if dest != "" {
				conf.Dest = dest
			}
			if conf.Dest == "" {
				return fmt.Errorf("destination must be specified in the config file or given by argument")
			}

			executor.Secret, err = conf.Secret.Open(c.Context)
			if err != nil {
				return fmt.Errorf("open secret store: %w", err)
			}

			num_errors := 0
			executor.NewHook = func(ctx context.Context, t task.Task) hook.Hook {
				return hook.Join(
					&countErrHook{n: &num_errors},
					hook.NewPrintHook(os.Stdout, t),
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

				executor.Execute(c.Context, task)
				i++
			})

			if num_errors > 0 {
				return cli.Exit("failed to bring some of things", 1)
			}

			return nil
		},
	}
}

type countErrHook struct {
	hook.NopHook
	n *int
}

func (h *countErrHook) OnError(err error) {
	*h.n++
}
