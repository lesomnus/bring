package cmd

import "github.com/urfave/cli/v2"

func NewApp() *cli.App {
	root := NewCmdBring()

	return &cli.App{
		Name:  "bring",
		Usage: "Bring things",
		Flags: root.Flags,
		Commands: []*cli.Command{
			NewCmdDigest(),
		},
		Action: root.Action,
	}
}
