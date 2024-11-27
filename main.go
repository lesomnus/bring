package main

import (
	"context"
	"errors"
	"log/slog"
	"os"

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

	app := cmd.NewApp(&l)
	app.ExitErrHandler = func(ctx context.Context, cmd *cli.Command, err error) {
		handleError(err)
	}
	if err := app.Run(context.Background(), os.Args); err != nil {
		handleError(err)
	}
}
