package main

import (
	"github.com/denudge/auto-updater/database"
	"github.com/denudge/auto-updater/migrations"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
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
			{
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
				},
			},
		},
	}
}
