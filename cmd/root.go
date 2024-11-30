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

func MewCmdRoot(logger **slog.Logger) *cli.Command {
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

		Description: "Bring files from the various source into the directory declaratively with integrity.",
		UsageText: `bring [GLOBAL OPTIONS] [BRING OPTIONS] [INVENTORY_FILE:-"./things.yaml"] --into <DESTINATION_DIR>
bring [GLOBAL OPTIONS] COMMAND [COMMAND OPTIONS]

INVENTORY_FILE:
   File that describes things to bring.

DESTINATION_DIR:
   Directory where things to be placed.`,

		Flags: flags,
		Commands: []*cli.Command{
			NewCmdDigest(),
			NewCmdVersion(),
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			conf_path := ".bring.yaml"
			conf_path_opt := cmd.HasName("conf")
			if conf_path_opt {
				conf_path = cmd.String("conf")
			}

			c, err_read := config.FromPath(conf_path)
			if err_read != nil {
				c = config.New()
			}
			if cmd.Bool("verbose") {
				c.Log.Level = "info"
			}
			if v := cmd.String("log-level"); v != "" {
				c.Log.Level = v
			}
			if v := cmd.String("log-format"); v != "" {
				c.Log.Format = v
			}

			l := c.Log.Logger()
			*logger = l
			if err_read == nil {
				l.Info("read config from the file", slog.String("path", conf_path))
			} else {
				if !errors.Is(err_read, os.ErrNotExist) || conf_path_opt {
					// There is an error
					// - config on default path
					// - config on given path
					err := fmt.Errorf("read config: %w", err_read)
					return nil, err
				}

				// There is no config on default path, so use default config.
				l.Info("use default config")
			}

			ctx = config.Into(ctx, c)
			ctx = log.Into(ctx, l)
			return ctx, nil
		},
		Action: root.Action,

		CustomRootCommandHelpTemplate: root_help_template,
	}
}

var root_help_template = `NAME:
   {{ template "helpNameTemplate" . }}

USAGE:
   {{ wrap .UsageText  3 }}

BRING OPTIONS:
{{- range .VisibleFlagCategories -}}
{{ if eq .Name "Bring" -}}
{{ range $i, $e := .Flags }}   {{ $e }}
{{ end -}}
{{ end }}
{{ end -}}
{{ if not .HideVersion -}}
VERSION:
   {{ .Version }}

{{ end -}}
{{ if .Description -}}
DESCRIPTION:
   {{ template "descriptionTemplate" . }}

{{ end -}}
{{ if .VisibleCommands -}}
COMMANDS:
{{- template "visibleCommandCategoryTemplate" . }}

{{ end -}}
GLOBAL OPTIONS:
{{- template "visibleFlagTemplate" . }}
`
