package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/lesomnus/bring/bringer"
	"github.com/lesomnus/bring/thing"
	"github.com/opencontainers/go-digest"
	"github.com/urfave/cli/v2"
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
				Action: func(ctx *cli.Context, s string) error {
					algorithm = digest.Algorithm(s)
					if algorithm.Available() {
						return nil
					}

					return fmt.Errorf("unknown hash algorithm; must be one of sha256, sha384 or sha512")
				},
			},
		},
		Action: func(c *cli.Context) error {
			target := c.Args().Get(0)
			if !strings.Contains(target, "://") {
				target = "file://" + target
			}

			url, err := url.Parse(target)
			if err != nil {
				return fmt.Errorf("parse target: %w", err)
			}

			b, err := bringer.FromUrl(url)
			if err != nil {
				return fmt.Errorf("get bringer: %w", err)
			}

			t := thing.Thing{Url: url}
			r, err := b.Bring(c.Context, t)
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
