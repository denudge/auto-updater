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

	args := os.Args

	/*
		args := []string{
			"bin/catalog",
			"app",
			"set-default-groups",
			"--vendor",
			"Foo",
			"--product",
			"Bar",
			"--default-group",
			"Betatester",
		}
	*/

	if err := cmd.Run(args); err != nil {
		log.Fatal(err)
	}
}
