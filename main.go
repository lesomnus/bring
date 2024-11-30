package main

import (
	"context"
	_ "crypto/sha256"
	_ "crypto/sha512"
	"errors"
	"log/slog"
	"os"

	_ "github.com/lesomnus/bring/bringer/file"
	_ "github.com/lesomnus/bring/bringer/http"
	_ "github.com/lesomnus/bring/bringer/smb"
	"github.com/lesomnus/bring/cmd"
	"github.com/urfave/cli/v3"
)

func main() {
	var l *slog.Logger
	handleError := func(err error) {
		if l != nil {
			l.Error(err.Error())
		}

		var e cli.ExitCoder
		if errors.As(err, &e) {
			os.Exit(e.ExitCode())
		} else {
			os.Exit(1)
		}
	}

	app := cmd.MewCmdRoot(&l)
	app.ExitErrHandler = func(ctx context.Context, cmd *cli.Command, err error) {
		handleError(err)
	}
	if err := app.Run(context.Background(), os.Args); err != nil {
		handleError(err)
	}
}
