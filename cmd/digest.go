package cmd

import (
	"context"
	"fmt"
	"net/url"

	"github.com/lesomnus/bring/bringer"
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
				return fmt.Errorf("parse target: %w", err)
			} else {
				u = *v
			}

			b, err := bringer.FromUrl(u)
			if err != nil {
				return fmt.Errorf("get bringer: %w", err)
			}

			t := thing.Thing{Url: u}
			r, err := b.Bring(ctx, t)
			if err != nil {
				return fmt.Errorf("get reader: %w", err)
			}
			defer r.Close()

			d, err := algorithm.FromReader(r)
			if err != nil {
				return fmt.Errorf("read: %w", err)
			}

			fmt.Println(d.String())
			return nil
		},
	}
}
