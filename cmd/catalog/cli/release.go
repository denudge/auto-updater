package cli

import (
	"fmt"
	"github.com/denudge/auto-updater/catalog"
	"github.com/urfave/cli/v2"
	"time"
)

func (console *Console) createReleaseCommands() *cli.Command {
	return &cli.Command{
		Name:  "release",
		Usage: "release management",
		Subcommands: []*cli.Command{
			{
				Name:  "latest",
				Usage: "list recently published releases",
				Flags: append(createReleaseFilterFlags(), createLimitFlag(10)[0]),
				Action: func(c *cli.Context) error {
					limit := parseLimitFlag(c, 10)
					return console.app.ListLatestReleases(limit)
				},
			},
			{
				Name:  "list",
				Usage: "list specific releases",
				Flags: append(createReleaseFilterFlags(), createLimitFlag(0)[0]),
				Before: func(c *cli.Context) error {
					groups := c.StringSlice("group")
					if groups != nil && len(groups) > 0 {
						return checkGroupsInput(groups)
					}

					return nil
				},
				Action: func(c *cli.Context) error {
					limit := parseLimitFlag(c, 0)
					filter := parseReleaseFilterFlags(c)
					releases, err := console.app.Store.FetchReleases(filter, limit)
					if err != nil {
						return err
					}

					printReleases(releases)

					return nil
				},
			},
			{
				Name:  "show",
				Usage: "show release details",
				Flags: createReleaseFilterFlags(),
				Before: func(c *cli.Context) error {
					groups := c.StringSlice("group")
					if groups != nil && len(groups) > 0 {
						return checkGroupsInput(groups)
					}

					return nil
				},
				Action: func(c *cli.Context) error {
					filter := parseReleaseFilterFlags(c)
					releases, err := console.app.Store.FetchReleases(filter, 0)
					if err != nil {
						return err
					}

					if releases == nil || len(releases) < 1 {
						fmt.Println("No release found.")

						return nil
					}

					if releases != nil && len(releases) > 1 {
						fmt.Printf("Filter is ambiguous. Please set more specific filters.\n\n")
						fmt.Println("Candidates:")
						printReleases(releases)

						return nil
					}

					printReleaseDetails(releases[0])

					return nil
				},
			},
			{
				Name:  "publish",
				Usage: "publish new release",
				Flags: createReleaseFlags(),
				Before: func(c *cli.Context) error {
					err := checkArguments(c, "publish", []string{"version"})
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

					storedApp, err := console.app.Store.FindApp(release.App.Vendor, release.App.Product)
					if err != nil || storedApp == nil {
						fmt.Printf("Console \"%s\" not found. Please create the console first.\n", release.App.String())

						return nil
					}

					stored, err := console.app.Store.StoreRelease(release, false)
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
				Flags: append(createReleaseFilterFlags(), &cli.StringFlag{Name: "upgrade-target", Usage: "The desired upgrade target", Required: true}),
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
		&cli.StringFlag{Name: "vendor", Usage: "Vendor name", Required: true},
		&cli.StringFlag{Name: "product", Usage: "Product name", Required: true},
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
		&cli.StringFlag{Name: "vendor", Usage: "Vendor name", Required: true},
		&cli.StringFlag{Name: "product", Usage: "Product name", Required: true},
		&cli.StringFlag{Name: "min-version", Usage: "Minimal version in semantic versioning scheme"},
		&cli.StringFlag{Name: "after-version", Usage: "Minimal excluded version in semantic versioning scheme"},
		&cli.StringFlag{Name: "before-version", Usage: "Maximum excluded version in semantic versioning scheme"},
		&cli.StringFlag{Name: "max-version", Usage: "Maximum version in semantic versioning scheme"},
		&cli.StringFlag{Name: "name", Usage: "Product name (for printing)"},
		&cli.StringFlag{Name: "variant", Usage: "Variant (Pro, Free, ...)"},
		&cli.StringFlag{Name: "os", Usage: "Operating system (MacOS, darwin, linux, ...)"},
		&cli.StringFlag{Name: "arch", Usage: "Architecture (i686, ppc64, ...)"},
		&cli.StringFlag{Name: "alias", Usage: "Alias name for the release"},
		&cli.BoolFlag{Name: "with-unstable", Usage: "Include unstable releases", DefaultText: "false"},
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

func printReleases(releases []*catalog.Release) {
	for _, release := range releases {
		fmt.Printf("%s\n", release)
	}
}

func printReleaseDetails(release *catalog.Release) {
	fmt.Printf("Vendor: %s\n", release.App.Vendor)
	fmt.Printf("Product: %s\n", release.App.Product)
	fmt.Printf("Variant: %s\n", release.Variant)
	fmt.Printf("Version: %s\n", release.Version)
	fmt.Printf("Alias: %s\n", release.Alias)
	fmt.Printf("OS: %s\n", release.OS)
	fmt.Printf("Arch: %s\n", release.Arch)
	fmt.Printf("Published: %s\n", release.Date.Format(time.RFC1123))
	fmt.Printf("Unstable: %v\n", release.Unstable)
	fmt.Printf("Groups: %s\n", catalog.FormatGroups(release.Groups))
	fmt.Printf("Upgrade Target: %s\n", release.UpgradeTarget)
	fmt.Printf("Criticality: %v\n", release.ShouldUpgrade)
	fmt.Printf("Format: %s\n", release.Format)
	fmt.Printf("Link: %s\n", release.Link)
}
