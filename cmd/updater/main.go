package main

import (
	"github.com/denudge/auto-updater/cmd/updater/app"
	"github.com/denudge/auto-updater/cmd/updater/cli"
	"log"
	"os"
)

func main() {
	// Users can override (re-init) these values with CLI parameters
	updaterApp := app.NewUpdater("updater.cfg", "")

	client := cli.NewConsole(updaterApp)

	args := os.Args

	if err := client.Cli.Run(args); err != nil {
		log.Fatal(err)
	}
}
