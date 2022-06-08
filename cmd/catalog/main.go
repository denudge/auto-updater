package main

import (
	"context"
	"github.com/denudge/auto-updater/cmd/catalog/api"
	"github.com/denudge/auto-updater/cmd/catalog/app"
	"github.com/denudge/auto-updater/cmd/catalog/cli"
	"github.com/denudge/auto-updater/config"
	"github.com/denudge/auto-updater/database"
	"log"
	"os"
)

func main() {
	db := database.Connect(config.Get("POSTGRES_DSN"))
	defer database.Close(db)

	// Create app
	cApp := app.NewApp(db, context.Background())

	// Create API
	cApi := api.NewApi(cApp)

	// Create CLI
	console := cli.NewConsole(cApp, cApi)

	args := os.Args

	if err := console.Cli.Run(args); err != nil {
		log.Fatal(err)
	}
}
