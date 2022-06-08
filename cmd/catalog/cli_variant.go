package main

import (
	"fmt"
	"github.com/denudge/auto-updater/catalog"
	"github.com/urfave/cli/v2"
	"time"
)

func (app *App) createVariantCommands() *cli.Command {
	return &cli.Command{
		Name:  "variant",
		Usage: "variant management",
		Subcommands: []*cli.Command{
			{
				Name:  "create",
				Usage: "create a new variant",
				Flags: createFullVariantFlags(),
				Before: func(c *cli.Context) error {
					defaultGroups := c.StringSlice("default-group")
					if defaultGroups != nil && len(defaultGroups) > 0 {
						return checkGroupsInput(defaultGroups)
					}

					return nil
				},
				Action: func(c *cli.Context) error {
					v := parseVariantFlags(c)
					v.Created = time.Now()

					stored, err := app.store.StoreVariant(v, false)
					if err != nil {
						return err
					}

					// time.Time.Before() cannot be used because the database might drop fractional seconds
					if stored.Created.Unix() < v.Created.Unix() {
						fmt.Printf("Variant \"%s\" has already been there.\n", stored)
					} else {
						fmt.Printf("Variant \"%s\" have been created.\n", stored)
					}

					return nil
				},
			},
			{
				Name:  "list",
				Usage: "list variants",
				Flags: append(createVariantFlags(), createLimitFlag(0)[0]),
				Action: func(c *cli.Context) error {
					v := parseVariantFlags(c)

					limit := parseLimitFlag(c, 0)

					return app.listAppVariants(v.App.Vendor, v.App.Product, v.Name, limit)
				},
			},
		},
	}
}

func (app *App) listAppVariants(vendor string, product string, name string, limit int) error {

	filter := catalog.VariantFilter{
		Vendor:  vendor,
		Product: product,
		Name:    name,
	}

	variants, err := app.store.ListVariants(filter, limit)
	if err != nil {
		return err
	}

	for _, variant := range variants {
		fmt.Printf("%s\n", variant.Name)
	}

	return nil
}

func createVariantFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "vendor", Usage: "Vendor name", Required: true},
		&cli.StringFlag{Name: "product", Usage: "Product name", Required: true},
		&cli.StringFlag{Name: "name", Usage: "Variant name"},
	}
}

func createFullVariantFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "vendor", Usage: "Vendor name", Required: true},
		&cli.StringFlag{Name: "product", Usage: "Product name", Required: true},
		&cli.StringFlag{Name: "name", Usage: "Variant name", Required: true},
		&cli.BoolFlag{Name: "inactive", Usage: "", DefaultText: "false"},
		&cli.BoolFlag{Name: "allow-register", Usage: "If clients can register", DefaultText: "false"},
		&cli.BoolFlag{Name: "locked", Usage: "", DefaultText: "false"},
		&cli.StringFlag{Name: "upgrade-target", Usage: "Optional: upgrade target for the app variant"},
		&cli.StringSliceFlag{Name: "default-group", Usage: "Optional: Default group(s). Specify none to use the app's default groups."},
	}
}

func parseVariantFlags(c *cli.Context) *catalog.Variant {
	return &catalog.Variant{
		App: &catalog.App{
			Vendor:  c.String("vendor"),
			Product: c.String("product"),
		},
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
