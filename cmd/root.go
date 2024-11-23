package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/lesomnus/bring/config"
	"github.com/lesomnus/bring/log"
	"github.com/urfave/cli/v2"
)

func NewApp() *cli.App {
	root := NewCmdBring()
	flags := []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   `Set log level to "info"`,
		},
		&cli.StringFlag{
			Name:  "log-level",
			Usage: `Set log level ["error" | "warn" | "info" | "debug"]`,
		},
		&cli.StringFlag{
			Name:  "log-format",
			Usage: `Set log format ["text" | "json" | "simple"]`,
		},
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Path to config file",

			Value: "things.yaml",

			TakesFile: true,
		},
	}

	flags = append(flags, root.Flags...)

	return &cli.App{
		Name:  "bring",
		Usage: "Bring things.",

		UsageText: `bring [GLOBAL OPTIONS] INVENTORY [DESTINATION]
bring [GLOBAL OPTIONS] COMMAND [COMMAND OPTIONS]

Example:

	bring things.yaml ./inventory/

		Bring things and use config described in the "thing.yaml" but
		the destination is overridden by "./inventory/".

	bring --conf conf.yaml things.yaml

		Bring things described in the "things.yaml" but use the config
		in the "conf.yaml".`,

		Description: "Bring files from the various source into the directory declaratively with integrity.",

		Flags: flags,
		Commands: []*cli.Command{
			NewCmdDigest(),
		},
		Before: func(ctx *cli.Context) error {
			conf_path := "things.yaml"
			if v := ctx.String("conf"); v != "" {
				conf_path = v
			}

			c, err := config.LoadFromFilepath(conf_path)
			if err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("load config: %w", err)
				}

				c = config.New()
			}
			if ctx.Bool("verbose") {
				c.Log.Level = "info"
			}
			if v := ctx.String("log-level"); v != "" {
				c.Log.Level = v
			}
			if v := ctx.String("log-format"); v != "" {
				c.Log.Format = v
			}

			l := c.Log.Logger()
			if err != nil {
				l.Info("use default config")
			} else {
				l.Info("load config from the file", slog.String("path", conf_path))
			}

			ctx.Context = config.Into(ctx.Context, c)
			ctx.Context = log.Into(ctx.Context, l)
			return nil
		},
		Action: root.Action,

		ExitErrHandler: func(c *cli.Context, err error) {
			l := log.From(c.Context)

			var e cli.ExitCoder
			if errors.As(err, &e) {
				l.Error(e.Error())
				os.Exit(e.ExitCode())
			} else if err != nil {
				l.Error(err.Error())
				os.Exit(1)
			}
		},
	}
}
