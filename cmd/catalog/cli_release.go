package main

import (
	"fmt"
	"github.com/denudge/auto-updater/catalog"
	"github.com/urfave/cli/v2"
	"time"
)

func (app *App) createReleaseCommands() *cli.Command {
	return &cli.Command{
		Name:  "release",
		Usage: "release management",
		Subcommands: []*cli.Command{
			{
				Name:  "latest",
				Usage: "list recently published releases",
				Flags: append(createReleaseFilterFlags(), createLimitFlag(10)[0]),
				Before: func(c *cli.Context) error {
					return checkArguments(c, "latest", []string{"vendor", "product"})
				},
				Action: func(c *cli.Context) error {
					limit := parseLimitFlag(c, 10)
					return app.ListLatestReleases(limit)
				},
			},
			{
				Name:  "list",
				Usage: "list specific releases",
				Flags: append(createReleaseFilterFlags(), createLimitFlag(0)[0]),
				Before: func(c *cli.Context) error {
					err := checkArguments(c, "list", []string{"vendor", "product"})
					if err != nil {
						return err
					}

					groups := c.StringSlice("group")
					if groups != nil && len(groups) > 0 {
						return checkGroupsInput(groups)
					}

					return nil
				},
				Action: func(c *cli.Context) error {
					limit := parseLimitFlag(c, 0)
					filter := parseReleaseFilterFlags(c)
					releases, err := app.store.FetchReleases(filter, limit)
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
					err := checkArguments(c, "publish", []string{"vendor", "product", "version"})
					if err != nil {
						return err
					}

					groups := c.StringSlice("group")
					if groups != nil && len(groups) > 0 {
						return checkGroupsInput(groups)
					}

					return nil
				},
				Action: func(c *cli.Context) error {
					release := parseReleaseFlags(c)

					storedApp, err := app.store.FindApp(release.App.Vendor, release.App.Product)
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
					} else {
						fmt.Printf("Release has been published: %s\n", stored)
					}

					return nil
				},
			},
			{
				Name:  "set-upgrade-target",
				Usage: "Set the upgrade target",
				Flags: append(createReleaseFilterFlags(), &cli.StringFlag{Name: "upgrade-target", Usage: "The desired upgrade target"}),
				Before: func(c *cli.Context) error {
					err := checkArguments(c, "set-upgrade-target", []string{"vendor", "product", "upgrade-target"})
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
			{
				Name:  "set-groups",
				Usage: "Bind releases to a new set of groups or make it public",
				Flags: func() []cli.Flag {
					flags := createReleaseFilterFlags()
					return append(flags[:len(flags)-1], &cli.StringSliceFlag{Name: "group", Usage: "Group(s). Use a single \"public\" group to specify the public group.", Required: true})
				}(),
				Before: func(c *cli.Context) error {
					err := checkArguments(c, "set-groups", []string{"vendor", "product"})
					if err != nil {
						return err
					}

					return checkGroupsInput(c.StringSlice("group"))
				},
				Action: func(c *cli.Context) error {
					groups := c.StringSlice("group")

					fmt.Printf("Setting the groups to %v\n", groups)

					return nil
				},
			},
		},
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
		&cli.StringSliceFlag{Name: "group", Usage: "Optional: Group(s). Use a single \"public\" group to specify the public group."},
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
		Groups:        c.StringSlice("group"),
	}
}

func createReleaseFilterFlags() []cli.Flag {
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
		&cli.StringSliceFlag{Name: "group", Usage: "Optional: Group(s). Use a single \"public\" group to specify the public group."},
	}
}

func parseReleaseFilterFlags(c *cli.Context) catalog.Filter {
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
		Groups:        c.StringSlice("group"),
	}

	filter.CompleteVersions()

	return filter
}
