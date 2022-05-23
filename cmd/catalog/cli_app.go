package main

import (
	"fmt"
	"github.com/denudge/auto-updater/catalog"
	"github.com/urfave/cli/v2"
)

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
					return checkArguments(c, "create", []string{"vendor", "product"})
				},
				Action: func(c *cli.Context) error {
					a := parseAppFlags(c)

					stored, err := app.store.StoreApp(a, false)
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

					limit := parseLimitFlag(c, 0)
					return app.ListApps(limit)
				},
			},
			{
				Name:  "set-default-groups",
				Usage: "Sets the default groups for an app",
				Flags: append(createAppFlags(), &cli.StringSliceFlag{Name: "default-group", Usage: "Default group(s). Specify a single \"public\" group to unlink special groups."}),
				Before: func(c *cli.Context) error {
					return checkArguments(c, "create", []string{"vendor", "product", "default-group"})
				},
				Action: func(c *cli.Context) error {
					a := parseAppFlags(c)

					if a.DefaultGroups != nil {
						err := checkGroupsInput(a.DefaultGroups)
						if err != nil {
							return err
						}
					}

					a, err := app.store.SetAppDefaultGroups(a)
					if err != nil {
						return err
					}

					app.printApp(a)

					return nil
				},
			},
		},
	}
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

func (app *App) ListApps(limit int) error {
	apps, err := app.store.ListApps(limit)
	if err != nil {
		return err
	}

	for _, a := range apps {
		app.printApp(a)
	}

	return nil
}

func (app *App) printApp(a *catalog.App) {
	groups := ""
	if a.DefaultGroups != nil && len(a.DefaultGroups) > 0 {
		defaultGroups := catalog.FormatGroups(a.DefaultGroups)
		groups = defaultGroups[1 : len(defaultGroups)-1]
	} else {
		groups = "(public)"
	}

	fmt.Printf("%s, default groups: %s\n", a, groups)
}
