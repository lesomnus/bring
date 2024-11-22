package cmd

import (
	"fmt"

	"github.com/lesomnus/bring/config"
	"github.com/urfave/cli/v2"
)

func NewCmdBring() *cli.Command {
	executor := &Executor{
		DryRun:  false,
		NewHook: NewStdIoPrinterHook,
	}

	return &cli.Command{
		Name:  "bring",
		Usage: "Bring things",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "dry-run",
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

			executor.Context = ExecuteContext{
				Context: c.Context,

				N: conf.Things.Len(),
				I: 0,
			}
			conf.Things.Walk(conf.Dest, executor.Execute)

			return nil
		},
	}
}
