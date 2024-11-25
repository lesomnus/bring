package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

type buildInfo struct {
	Version   string
	TimeBuild string
	GitRev    string
	GitDirty  bool
}

//go:generate go run ./gen/version
var _buildInfo = buildInfo{
	Version:   "v0.0.0-edge",
	TimeBuild: time.Now().Format(time.RFC3339),
	GitRev:    "0000000000000000000000000000000000000000",
	GitDirty:  true,
}

func NewCmdVersion() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Print the program information",

		Action: func(c *cli.Context) error {
			b := strings.Builder{}
			b.WriteString(fmt.Sprintf("BRING_VERSION=%s\n", _buildInfo.Version))
			b.WriteString(fmt.Sprintf("BRING_TIME_BUILD=%s\n", _buildInfo.TimeBuild))
			b.WriteString(fmt.Sprintf("BUILD_GIT_REV=%s", _buildInfo.GitRev))
			if _buildInfo.GitDirty {
				b.WriteString("-dirty")
			}
			b.WriteString("\n")

			fmt.Print(b.String())
			return nil
		},
	}
}
