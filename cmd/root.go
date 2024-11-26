package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/lesomnus/bring/config"
	"github.com/lesomnus/bring/log"
	"github.com/urfave/cli/v3"
)

func NewApp(l **slog.Logger) *cli.Command {
	root := NewCmdBring()
	flags := []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   `Set log level to "info"`,
		},
		&cli.StringFlag{
			Name:  "log-level",
			Usage: `"error" | "warn" | "info" | "debug"`,
		},
		&cli.StringFlag{
			Name:  "log-format",
			Usage: `"text" | "json" | "simple"`,
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

	return &cli.Command{
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
			NewCmdVersion(),
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			conf_path := "things.yaml"
			conf_path_opt := cmd.HasName("conf")
			if conf_path_opt {
				conf_path = cmd.String("conf")
			}

			conf, err := config.LoadFromFilepath(conf_path)
			if err != nil {
				conf = config.New()
				if !errors.Is(err, os.ErrNotExist) || conf_path_opt {
					// There is an error
					// - config on default path
					// - config on explicitly given path
					err := fmt.Errorf("load config: %w", err)
					*l = conf.Log.Logger()
					return nil, err
				}

				// There is no config on default path, so use default config.
				defer func() {
					(*l).Info("use default config")
				}()
			} else {
				defer func() {
					(*l).Info("load config from the file", slog.String("path", conf_path))
				}()
			}
			if cmd.Bool("verbose") {
				conf.Log.Level = "info"
			}
			if v := cmd.String("log-level"); v != "" {
				conf.Log.Level = v
			}
			if v := cmd.String("log-format"); v != "" {
				conf.Log.Format = v
			}

			*l = conf.Log.Logger()

			ctx = config.Into(ctx, conf)
			ctx = log.Into(ctx, *l)
			return ctx, nil
		},
		Action: root.Action,
	}
}
