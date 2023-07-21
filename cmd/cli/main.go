package main

import (
	"fmt"
	"os"

	"github.com/owlint/lokal/internal/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Lokal CLI",
		Usage: "Easily run your kubernetes applications in local.",
		Commands: []*cli.Command{
			commands.RunCommand,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
