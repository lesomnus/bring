package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

type buildInfo struct {
	Version   string
	BuildTime string
	GitHash   string
	GitDirty  bool
}

//go:generate go run ./gen/version
var _buildInfo = buildInfo{
	Version:   "v0.0.0-edge",
	BuildTime: time.Now().Format(time.RFC3339),
	GitHash:   "0000000000000000000000000000000000000000",
	GitDirty:  true,
}

func NewCmdVersion() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Print the program information",

		Action: func(c *cli.Context) error {
			b := strings.Builder{}
			b.WriteString(fmt.Sprintf("Version %s\n", _buildInfo.Version))
			b.WriteString(fmt.Sprintf("Built at %s\n", _buildInfo.BuildTime))
			b.WriteString(fmt.Sprintf("Revision %s", _buildInfo.GitHash))
			if _buildInfo.GitDirty {
				b.WriteString(" dirty")
			}
			b.WriteString("\n")

			fmt.Print(b.String())
			return nil
		},
	}
}
