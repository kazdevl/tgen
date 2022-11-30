package main

import (
	"log"
	"os"

	"github.com/kazdevl/tgen/cmd/tgen/subcmd"
	"github.com/urfave/cli/v2"
)

const (
	version = "1.0.0"
)

func main() {
	app := &cli.App{
		Name:     "tgen",
		Usage:    "make go test code with mock",
		Version:  version,
		Commands: subcmd.ProvideSubCommands(),
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
