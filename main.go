package main

import (
	"context"
	"os"

	"github.com/lesomnus/bring/cmd"
)

func main() {
	app := cmd.NewApp()
	app.Run(context.Background(), os.Args)
}
