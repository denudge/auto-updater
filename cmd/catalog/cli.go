package main

import (
	"fmt"
	"github.com/denudge/auto-updater/catalog"
	"github.com/denudge/auto-updater/database"
	"github.com/denudge/auto-updater/migrations"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
	"time"
)

func NewCli(app *App, api *Api) *cli.App {
	return &cli.App{
		Name: "catalog",
		Commands: []*cli.Command{
			// A bunch of database (migration) related commands
			database.NewCommand(migrate.NewMigrator(app.db, migrations.Migrations)),
			{
				Name:  "serve",
				Usage: "run HTTP catalog API server",
				Action: func(c *cli.Context) error {
					api.Serve()
					return nil
				},
			},
			app.createReleaseCommands(),
		},
	}
}

func (app *App) createReleaseCommands() *cli.Command {
	return &cli.Command{
		Name:  "release",
		Usage: "release management",
		Subcommands: []*cli.Command{
			{
				Name:  "latest",
				Usage: "List recently published releases",
				Action: func(c *cli.Context) error {
					return app.ListLatestReleases()
				},
			},
			{
				Name:  "publish",
				Usage: "Publish new release",
				Action: func(c *cli.Context) error {

					// TODO: Read values from arguments
					release := &catalog.Release{
						Vendor:  "Foo",
						Product: "Bar",
						Variant: "Pro",
						OS:      "",
						Arch:    "",
						Version: "1.1.0",
						Date:    time.Now(),
					}

					stored, err := app.store.Store(release, false)
					if err != nil {
						return err
					}

					if stored.Date.Before(release.Date) {
						fmt.Println("Release has already been there.")
					}

					return nil
				},
			},
		},
	}
}
