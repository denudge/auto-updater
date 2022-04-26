package main

import (
	"context"
	"github.com/denudge/auto-updater/config"
	"github.com/denudge/auto-updater/database"
	"log"
	"os"
)

func main() {
	db := database.Connect(config.Get("POSTGRES_DSN"))
	defer database.Close(db)

	// Create app
	app := NewApp(db, context.Background())

	// Create API
	api := NewApi(app)

	// Create CLI
	cmd := NewCli(app, api)

	if err := cmd.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
