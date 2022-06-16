package main

import (
	"github.com/denudge/auto-updater/cmd/updater/app"
	"github.com/denudge/auto-updater/cmd/updater/cli"
	"log"
	"os"
)

const DefaultConfigFileName = "updater.cfg"

func main() {
	// Users can override (re-init) these values with CLI parameters
	// or they are read from a given config file
	updaterApp := app.NewUpdater(DefaultConfigFileName, "")

	client := cli.NewConsole(updaterApp)

	args := os.Args

	if err := client.Cli.Run(args); err != nil {
		log.Fatal(err)
	}
}
