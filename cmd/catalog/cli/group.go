package cli

import (
	"errors"
	"fmt"
	"github.com/denudge/auto-updater/catalog"
	"github.com/urfave/cli/v2"
	"time"
)

func (console *Console) createGroupCommands() *cli.Command {
	return &cli.Command{
		Name:  "group",
		Usage: "group management",
		Subcommands: []*cli.Command{
			{
				Name:  "create",
				Usage: "create a new group",
				Flags: createGroupFlags(),
				Before: func(c *cli.Context) error {
					return checkArguments(c, "create", []string{"name"})
				},
				Action: func(c *cli.Context) error {
					g := parseGroupFlags(c)
					g.Created = time.Now()

					stored, err := console.app.Store.StoreGroup(g, false)
					if err != nil {
						return err
					}

					// time.Time.Before() cannot be used because the database might drop fractional seconds
					if stored.Created.Unix() < g.Created.Unix() {
						fmt.Println("Group has already been there.")
					} else {
						fmt.Printf("Group \"%s\" have been created.\n", stored.Name)
					}

					return nil
				},
			},
			{
				Name:  "list",
				Usage: "list groups",
				Flags: append(createGroupFlags(), createLimitFlag(0)[0]),
				Action: func(c *cli.Context) error {
					g := parseGroupFlags(c)

					limit := parseLimitFlag(c, 0)

					return console.listAppGroups(g.App.Vendor, g.App.Product, g.Name, limit)
				},
			},
		},
	}
}

func (console *Console) listAppGroups(vendor string, product string, name string, limit int) error {

	filter := catalog.GroupFilter{
		Vendor:  vendor,
		Product: product,
		Name:    name,
	}

	groups, err := console.app.Store.ListGroups(filter, limit)
	if err != nil {
		return err
	}

	for _, group := range groups {
		fmt.Printf("%s\n", group.Name)
	}

	return nil
}

func createGroupFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "vendor", Usage: "Vendor name", Required: true},
		&cli.StringFlag{Name: "product", Usage: "Product name", Required: true},
		&cli.StringFlag{Name: "name", Usage: "Group name"},
	}
}

func parseGroupFlags(c *cli.Context) *catalog.Group {
	return &catalog.Group{
		App: &catalog.App{
			Vendor:  c.String("vendor"),
			Product: c.String("product"),
		},
		Name:    c.String("name"),
		Created: time.Now(),
		Updated: time.Now(),
	}
}

func checkGroupsInput(strs []string) error {
	hasPublic := false
	hasOther := false

	if strs == nil || len(strs) < 1 {
		return errors.New("no groups given")
	}

	for _, str := range strs {
		if str == "public" {
			hasPublic = true
		} else {
			hasOther = true
		}
	}

	if hasPublic && hasOther {
		return errors.New("public and groups cannot be mixed")
	}

	return nil
}
