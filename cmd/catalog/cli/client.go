package cli

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func (console *Console) createClientCommands() *cli.Command {
	return &cli.Command{
		Name:  "client",
		Usage: "client management",
		Subcommands: []*cli.Command{
			{
				Name:  "register",
				Usage: "register an console client",
				Flags: append(createClientAppFlags(), &cli.StringSliceFlag{Name: "group", Usage: "Client group(s). Use none to put the client into the public group."}),
				Action: func(c *cli.Context) error {
					a := parseAppFlags(c)
					variant := c.String("variant")

					groups := c.StringSlice("group")
					if groups != nil && len(groups) > 0 {
						err := checkGroupsInput(groups)
						if err != nil {
							return err
						}
					}

					stored, err := console.app.Store.RegisterClient(a, variant, c.StringSlice("group"))
					if err != nil {
						return err
					}

					// time.Time.Before() cannot be used because the database might drop fractional seconds
					if stored.Created.Unix() < a.Created.Unix() {
						fmt.Println("Client has already been there.")
					} else {
						fmt.Printf("Client registered: %s\n", stored.Uuid)
					}

					return nil
				},
			},
		},
	}
}

func createClientAppFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "vendor", Usage: "Vendor name", Required: true},
		&cli.StringFlag{Name: "product", Usage: "Product name", Required: true},
		&cli.StringFlag{Name: "variant", Usage: "The variant of the product", Required: true},
	}
}
