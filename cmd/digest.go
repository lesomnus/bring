package cmd

import (
	"context"
	"fmt"
	"net/url"

	"github.com/lesomnus/bring/config"
	"github.com/lesomnus/bring/internal/hook"
	"github.com/lesomnus/bring/internal/hooks"
	"github.com/lesomnus/bring/internal/task"
	"github.com/lesomnus/bring/log"
	"github.com/lesomnus/bring/thing"
	"github.com/opencontainers/go-digest"
	"github.com/urfave/cli/v3"
)

func NewCmdDigest() *cli.Command {
	algorithm := digest.SHA256

	return &cli.Command{
		Name:  "digest",
		Usage: "Calculate digest of a thing from a file or an URL",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "with",
				Value: algorithm.String(),
				Usage: "Algorithm to digest",
				Action: func(ctx context.Context, cmd *cli.Command, s string) error {
					algorithm = digest.Algorithm(s)
					if algorithm.Available() {
						return nil
					}

					return fmt.Errorf("unknown hash algorithm; must be one of sha256, sha384 or sha512")
				},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			target := ""
			switch cmd.NArg() {
			case 0:
				return fmt.Errorf("resource URL must be given")
			case 1:
				target = cmd.Args().Get(0)
			default:
				return fmt.Errorf("expected exactly 1 argument")
			}

			var u url.URL
			if v, err := url.Parse(target); err != nil {
				return fmt.Errorf("parse resource URL: %w", err)
			} else {
				u = *v
			}

			c := config.From(ctx)
			l := log.From(ctx)
			executor := &executor{
				DryRun: true,
				NewHook: func(ctx context.Context, t task.Task) hook.Hook {
					return &hooks.LogHook{T: t, L: l}
				},
			}
			if err := c.Secret.OpenTo(ctx, u, &executor.Secret); err != nil {
				return err
			}

			ctx, cancel := c.Each.ApplyBringTimeout(ctx)
			defer cancel()

			r, err := executor.Execute(ctx, task.Task{
				Thing: thing.Thing{Url: u},

				BringConfig: c.Each,
			})
			if err != nil {
				return err
			} else {
				defer r.Close()
			}

			d, err := algorithm.FromReader(r)
			if err != nil {
				return fmt.Errorf("read: %w", err)
			}

			fmt.Println(d.String())
			return nil
		},
	}
}
