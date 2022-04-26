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
				Name:  "list",
				Usage: "list specific releases",
				Flags: createFilterFields(),
				Before: func(c *cli.Context) error {
					return checkFilterArguments(c, "list")
				},
				Action: func(c *cli.Context) error {
					filter := parseFilter(c)
					releases, err := app.store.Fetch(filter)
					if err != nil {
						return err
					}

					for _, release := range releases {
						fmt.Printf("%s\n", release)
					}

					return nil
				},
			},
			{
				Name:  "publish",
				Usage: "publish new release",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "vendor", Usage: "Vendor name"},
					&cli.StringFlag{Name: "product", Usage: "Product name"},
					&cli.StringFlag{Name: "version", Usage: "Version in semantic versioning scheme"},
					&cli.StringFlag{Name: "name", Usage: "Optional: product name (for printing)"},
					&cli.StringFlag{Name: "variant", Usage: "Optional: variant (Pro, Free, ...)"},
					&cli.StringFlag{Name: "os", Usage: "Optional: operating system (MacOS, darwin, linux, ...)"},
					&cli.StringFlag{Name: "arch", Usage: "Optional: architecture (i686, ppc64, ...)"},
					&cli.StringFlag{Name: "description", Usage: "Optional: notes"},
					&cli.StringFlag{Name: "alias", Usage: "Optional: alias name for the release"},
					&cli.StringFlag{Name: "upgrade-target", Usage: "Optional: upgrade target for the release"},
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
						Vendor:        c.String("vendor"),
						Product:       c.String("product"),
						Variant:       c.String("variant"),
						OS:            c.String("osName"),
						Arch:          c.String("arch"),
						Version:       c.String("version"),
						Date:          time.Now(),
						Name:          c.String("name"),
						Description:   c.String("description"),
						Alias:         c.String("alias"),
						Unstable:      c.Bool("unstable"),
						UpgradeTarget: catalog.UpgradeTarget(c.String("upgrade-target")),
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

func createFilterFields() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "vendor", Usage: "Vendor name"},
		&cli.StringFlag{Name: "product", Usage: "Product name"},
		&cli.StringFlag{Name: "min-version", Usage: "Minimal version in semantic versioning scheme"},
		&cli.StringFlag{Name: "after-version", Usage: "Minimal excluded version in semantic versioning scheme"},
		&cli.StringFlag{Name: "before-version", Usage: "Maximum excluded version in semantic versioning scheme"},
		&cli.StringFlag{Name: "max-version", Usage: "Maximum version in semantic versioning scheme"},
		&cli.StringFlag{Name: "name", Usage: "Product name (for printing)"},
		&cli.StringFlag{Name: "variant", Usage: "Variant (Pro, Free, ...)"},
		&cli.StringFlag{Name: "os", Usage: "Operating system (MacOS, darwin, linux, ...)"},
		&cli.StringFlag{Name: "arch", Usage: "Architecture (i686, ppc64, ...)"},
		&cli.StringFlag{Name: "alias", Usage: "Alias name for the release"},
		&cli.BoolFlag{Name: "with-unstable", Usage: "Include unstable releases"},
	}
}

func parseFilter(c *cli.Context) catalog.Filter {
	return catalog.Filter{
		Vendor:        c.String("vendor"),
		Product:       c.String("product"),
		Variant:       c.String("variant"),
		MinVersion:    c.String("min-version"),
		AfterVersion:  c.String("after-version"),
		BeforeVersion: c.String("before-version"),
		MaxVersion:    c.String("max-version"),
		OS:            c.String("osName"),
		Arch:          c.String("arch"),
		Name:          c.String("name"),
		Alias:         c.String("alias"),
		WithUnstable:  c.Bool("with-unstable"),
	}
}

func checkFilterArguments(c *cli.Context, command string) error {
	// Check arguments
	if c.String("vendor") == "" || c.String("product") == "" {
		_ = cli.ShowCommandHelp(c, command)

		return fmt.Errorf("At least vendor and product must be specified.\n")
	}

	return nil
}
