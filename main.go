package main

import (
	"os"

	"github.com/lesomnus/bring/cmd"
)

func main() {
	app := cmd.NewApp()
	app.Run(os.Args)
}
