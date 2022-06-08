package main

import (
	"fmt"
	"github.com/denudge/auto-updater/cmd/updater/app"
	"github.com/denudge/auto-updater/cmd/updater/cli"
	"log"
	"os"
	"strings"
)

func main() {

	catalogUrl := os.Getenv("CATALOG_URL")
	if catalogUrl == "" || !strings.HasPrefix(catalogUrl, "http") {
		fmt.Println("env variable CATALOG_URL missing or does not start with \"http\"")

		return
	}

	updater := app.NewUpdater(catalogUrl)

	client := cli.NewConsole(updater)

	args := os.Args

	if err := client.Cli.Run(args); err != nil {
		log.Fatal(err)
	}
}
