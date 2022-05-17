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
			app.createAppCommands(),
			app.createGroupCommands(),
			app.createReleaseCommands(),
		},
	}
}

func (app *App) createAppCommands() *cli.Command {
	return &cli.Command{
		Name:  "app",
		Usage: "app management",
		Subcommands: []*cli.Command{
			{
				Name:  "create",
				Usage: "create a new app",
				Flags: createAppFlags(),
				Before: func(c *cli.Context) error {
					return checkAppArguments(c, "create", "")
				},
				Action: func(c *cli.Context) error {
					a := parseAppFlags(c)

					stored, err := app.store.CreateApp(a, false)
					if err != nil {
						return err
					}

					// time.Time.Before() cannot be used because the database might drop fractional seconds
					if stored.Created.Unix() < a.Created.Unix() {
						fmt.Println("App has already been there.")
					}

					return nil
				},
			},
			{
				Name:  "list",
				Usage: "list apps",
				Flags: append(createAppFlags(), createLimitFlag()[0]),
				Action: func(c *cli.Context) error {

					limit := parseLimitFlag(c, 10)
					return app.ListApps(limit)
				},
			},
		},
	}
}

func (app *App) createGroupCommands() *cli.Command {
	return &cli.Command{
		Name:  "group",
		Usage: "group management",
		Subcommands: []*cli.Command{
			{
				Name:  "create",
				Usage: "create a new group",
				Flags: createGroupFlags(),
				Before: func(c *cli.Context) error {
					return checkAppArguments(c, "create", "name")
				},
				Action: func(c *cli.Context) error {
					g := parseGroupFlags(c)
					g.Created = time.Now()

					stored, err := app.store.StoreGroup(g, false)
					if err != nil {
						return err
					}

					// time.Time.Before() cannot be used because the database might drop fractional seconds
					if stored.Created.Unix() < g.Created.Unix() {
						fmt.Println("Group has already been there.")
					}

					return nil
				},
			},
			{
				Name:  "list",
				Usage: "list groups",
				Flags: append(createGroupFlags(), createLimitFlag()[0]),
				Before: func(c *cli.Context) error {
					return checkAppArguments(c, "list", "")
				},
				Action: func(c *cli.Context) error {
					g := parseGroupFlags(c)

					limit := parseLimitFlag(c, 10)
					filter := catalog.GroupFilter{
						Vendor:  g.App.Vendor,
						Product: g.App.Product,
						Name:    g.Name,
					}
					groups, err := app.store.ListGroups(filter, limit)
					if err != nil {
						return err
					}

					for _, group := range groups {
						fmt.Printf("%s\n", group.Name)
					}

					return nil
				},
			},
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
				Flags: append(createFilterFlags(), createLimitFlag()[0]),
				Before: func(c *cli.Context) error {
					return checkAppArguments(c, "latest", "")
				},
				Action: func(c *cli.Context) error {
					limit := parseLimitFlag(c, 10)
					return app.ListLatestReleases(limit)
				},
			},
			{
				Name:  "list",
				Usage: "list specific releases",
				Flags: append(createFilterFlags(), createLimitFlag()[0]),
				Before: func(c *cli.Context) error {
					return checkAppArguments(c, "list", "")
				},
				Action: func(c *cli.Context) error {
					filter := parseFilterFlags(c)
					releases, err := app.store.FetchReleases(filter)
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
				Flags: createReleaseFlags(),
				Before: func(c *cli.Context) error {
					// Check arguments
					if c.String("vendor") == "" || c.String("product") == "" || c.String("version") == "" {
						_ = cli.ShowCommandHelp(c, "publish")

						return fmt.Errorf("At least vendor, product and version must be specified.\n")
					}

					return nil
				},
				Action: func(c *cli.Context) error {
					release := parseReleaseFlags(c)

					storedApp, err := app.store.FindApp(release.App.Vendor, release.App.Vendor)
					if err != nil || storedApp == nil {
						fmt.Printf("App \"%s\" not found. Please create the app first.\n", release.App.String())

						return nil
					}

					stored, err := app.store.StoreRelease(release, false)
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
			{
				Name:  "set-upgrade-target",
				Usage: "Set the upgrade target",
				Flags: append(createFilterFlags(), &cli.StringFlag{Name: "upgrade-target", Usage: "The desired upgrade target"}),
				Before: func(c *cli.Context) error {
					err := checkAppArguments(c, "set-upgrade-target", "upgrade-target")
					if err != nil {
						return err
					}

					return nil
				},
				Action: func(c *cli.Context) error {
					value := c.String("upgrade-target")

					if !catalog.UpgradeTarget(value).IsValid() {
						return fmt.Errorf("is not a valid upgrade target: %s\n", value)
					}

					fmt.Printf("Setting the upgrade target to %s\n", value)

					return nil
				},
			},
		},
	}
}

func createLimitFlag() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{Name: "limit", Usage: ""},
	}
}

func parseLimitFlag(c *cli.Context, defaultValue int) int {
	given := c.String("limit")
	if given == "" {
		return defaultValue
	}

	return c.Int("limit")
}

func createAppFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "vendor", Usage: "Vendor name"},
		&cli.StringFlag{Name: "product", Usage: "Product name"},
		&cli.StringFlag{Name: "name", Usage: "Product name (for printing)"},
		&cli.BoolFlag{Name: "active", Usage: ""},
		&cli.BoolFlag{Name: "locked", Usage: ""},
		&cli.StringFlag{Name: "upgrade-target", Usage: "Optional: upgrade target for the app"},
	}
}

func createGroupFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "vendor", Usage: "Vendor name"},
		&cli.StringFlag{Name: "product", Usage: "Product name"},
		&cli.StringFlag{Name: "name", Usage: "Group name"},
	}
}

func parseAppFlags(c *cli.Context) *catalog.App {
	return &catalog.App{
		Vendor:        c.String("vendor"),
		Product:       c.String("product"),
		Name:          c.String("name"),
		Active:        c.Bool("active"),
		Locked:        c.Bool("locked"),
		UpgradeTarget: catalog.UpgradeTarget(c.String("upgrade-target")),
	}
}

func parseGroupFlags(c *cli.Context) *catalog.Group {
	return &catalog.Group{
		App: &catalog.App{
			Vendor:  c.String("vendor"),
			Product: c.String("product"),
		},
		Name: c.String("name"),
	}
}

func createReleaseFlags() []cli.Flag {
	return []cli.Flag{
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
	}
}

func parseReleaseFlags(c *cli.Context) *catalog.Release {
	return &catalog.Release{
		App: &catalog.App{
			Vendor:  c.String("vendor"),
			Product: c.String("product"),
			Name:    c.String("name"),
		},
		Variant:       c.String("variant"),
		OS:            c.String("os"),
		Arch:          c.String("arch"),
		Version:       c.String("version"),
		Date:          time.Now(),
		Description:   c.String("description"),
		Alias:         c.String("alias"),
		Unstable:      c.Bool("unstable"),
		UpgradeTarget: catalog.UpgradeTarget(c.String("upgrade-target")),
	}
}

func createFilterFlags() []cli.Flag {
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

func parseFilterFlags(c *cli.Context) catalog.Filter {
	filter := catalog.Filter{
		Vendor:        c.String("vendor"),
		Product:       c.String("product"),
		Variant:       c.String("variant"),
		MinVersion:    c.String("min-version"),
		AfterVersion:  c.String("after-version"),
		BeforeVersion: c.String("before-version"),
		MaxVersion:    c.String("max-version"),
		OS:            c.String("os"),
		Arch:          c.String("arch"),
		Name:          c.String("name"),
		Alias:         c.String("alias"),
		WithUnstable:  c.Bool("with-unstable"),
	}

	filter.CompleteVersions()

	return filter
}

func checkAppArguments(c *cli.Context, command string, additionalField string) error {
	// Check arguments
	if c.String("vendor") == "" || c.String("product") == "" ||
		(additionalField != "" && c.String(additionalField) == "") {
		_ = cli.ShowCommandHelp(c, command)

		requiredFields := "vendor and product"
		if additionalField != "" {
			requiredFields = "vendor, product and " + additionalField
		}

		return fmt.Errorf("At least %s must be specified.\n", requiredFields)
	}

	return nil
}
