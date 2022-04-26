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
				Usage: "list recently published releases",
				Action: func(c *cli.Context) error {
					return app.ListLatestReleases()
				},
			},
			{
				Name:  "publish",
				Usage: "publish new release",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "vendor", Usage: "Vendor name"},
					&cli.StringFlag{Name: "product", Usage: "Product name"},
					&cli.StringFlag{Name: "variant", Usage: "Variant (Pro, Free, ...)"},
					&cli.StringFlag{Name: "version", Usage: "Version in semantic versioning scheme"},
					&cli.StringFlag{Name: "os", Usage: "Operating system (MacOS, darwin, linux, ...)"},
					&cli.StringFlag{Name: "arch", Usage: "Architecture (i686, ppc64, ...)"},
					&cli.BoolFlag{Name: "unstable", Usage: "Mark release as unstable"},
				},
				Before: func(c *cli.Context) error {
					// Check arguments
					if c.String("vendor") == "" || c.String("product") == "" || c.String("version") == "" {
						_ = cli.ShowCommandHelp(c, "publish")

						return fmt.Errorf("At least vendor, product and version must be specified.\n")
					}

					return nil
				},
				Action: func(c *cli.Context) error {
					release := &catalog.Release{
						Vendor:   c.String("vendor"),
						Product:  c.String("product"),
						Variant:  c.String("variant"),
						OS:       c.String("osName"),
						Arch:     c.String("arch"),
						Version:  c.String("version"),
						Date:     time.Now(),
						Unstable: c.Bool("unstable"),
					}

					stored, err := app.store.Store(release, false)
					if err != nil {
						return err
					}

					// time.Time.Before() cannot be used because the database might drop fractional seconds
					if stored.Date.Unix() < release.Date.Unix() {
						fmt.Println("Release has already been there.")
					}

					return nil
				},
			},
		},
	}
}
