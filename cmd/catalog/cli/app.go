package cli

import (
	"fmt"
	"github.com/denudge/auto-updater/catalog"
	"github.com/urfave/cli/v2"
	"time"
)

func (console *Console) createAppCommands() *cli.Command {
	return &cli.Command{
		Name:  "console",
		Usage: "console management",
		Subcommands: []*cli.Command{
			{
				Name:  "create",
				Usage: "create a new console",
				Flags: createFullAppFlags(),
				Action: func(c *cli.Context) error {
					a := parseAppFlags(c)

					stored, err := console.app.Store.StoreApp(a, false)
					if err != nil {
						return err
					}

					// time.Time.Before() cannot be used because the database might drop fractional seconds
					if stored.Created.Unix() < a.Created.Unix() {
						fmt.Println("Console has already been there.")
					}

					return nil
				},
			},
			{
				Name:  "list",
				Usage: "list apps",
				Flags: createLimitFlag(0),
				Action: func(c *cli.Context) error {

					limit := parseLimitFlag(c, 0)
					return console.ListApps(limit)
				},
			},
			{
				Name:  "set-default-groups",
				Usage: "Sets the default groups for an console",
				Flags: append(createMinAppFlags(), &cli.StringSliceFlag{Name: "default-group", Usage: "Default group(s). Specify a single \"public\" group to unlink special groups."}),
				Before: func(c *cli.Context) error {
					return checkArguments(c, "create", []string{"default-group"})
				},
				Action: func(c *cli.Context) error {
					a := parseAppFlags(c)

					if a.DefaultGroups != nil {
						err := checkGroupsInput(a.DefaultGroups)
						if err != nil {
							return err
						}
					}

					a, err := console.app.Store.SetAppDefaultGroups(a)
					if err != nil {
						return err
					}

					console.printApp(a)

					return nil
				},
			},
		},
	}
}

func createMinAppFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "vendor", Usage: "Vendor name", Required: true},
		&cli.StringFlag{Name: "product", Usage: "Product name", Required: true},
	}
}

func createFullAppFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "vendor", Usage: "Vendor name", Required: true},
		&cli.StringFlag{Name: "product", Usage: "Product name", Required: true},
		&cli.StringFlag{Name: "name", Usage: "Product name (for printing)"},
		&cli.BoolFlag{Name: "inactive", Usage: "", DefaultText: "false"},
		&cli.BoolFlag{Name: "allow-register", Usage: "If clients can register", DefaultText: "false"},
		&cli.BoolFlag{Name: "locked", Usage: "", DefaultText: "false"},
		&cli.StringFlag{Name: "upgrade-target", Usage: "Optional: upgrade target for the app"},
	}
}

func parseAppFlags(c *cli.Context) *catalog.App {
	return &catalog.App{
		Vendor:        c.String("vendor"),
		Product:       c.String("product"),
		Name:          c.String("name"),
		Active:        !c.Bool("inactive"),
		Locked:        c.Bool("locked"),
		AllowRegister: c.Bool("allow-register"),
		UpgradeTarget: catalog.UpgradeTarget(c.String("upgrade-target")),
		DefaultGroups: c.StringSlice("default-group"),
		Created:       time.Now(),
		Updated:       time.Now(),
	}
}

func (console *Console) ListApps(limit int) error {
	apps, err := console.app.Store.ListApps(limit)
	if err != nil {
		return err
	}

	for _, a := range apps {
		console.printApp(a)
	}

	return nil
}

func (console *Console) printApp(a *catalog.App) {
	groups := ""
	if a.DefaultGroups != nil && len(a.DefaultGroups) > 0 {
		defaultGroups := catalog.FormatGroups(a.DefaultGroups)
		groups = defaultGroups[1 : len(defaultGroups)-1]
	} else {
		groups = "(public)"
	}

	fmt.Printf("%s, default groups: %s\n", a, groups)
}
