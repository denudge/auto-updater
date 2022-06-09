package cli

import (
	"fmt"
	"github.com/denudge/auto-updater/cmd/updater/app"
	"github.com/denudge/auto-updater/updater"
	"github.com/urfave/cli/v2"
	"runtime"
)

type Console struct {
	Cli     *cli.App
	updater *app.Updater
}

func NewConsole(updater *app.Updater) *Console {
	return &Console{
		Cli:     newCli(updater),
		updater: updater,
	}
}

func newCli(updater *app.Updater) *cli.App {
	return &cli.App{
		Name: "updater",
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "initialize the updater client",
				Flags: createInitializationFlags(),
				Action: func(c *cli.Context) error {
					// TODO: Check if config file already exists (and if, return error)

					state := parseInitializationFlags(c)

					if state.ClientId == "" {
						// We need to register
						response, err := updater.Client.RegisterClient(state.Vendor, state.Product, state.Variant)
						if err != nil {
							return err
						}

						if response == nil || response.ClientId == "" {
							return fmt.Errorf("could not register client (unknown reason)")
						}

						state.ClientId = response.ClientId

						fmt.Printf("Successfully registered with client ID %s\n", state.ClientId)
					}

					return nil
				},
			},
		},
	}
}

type ClientState struct {
	ClientId string
	Vendor   string
	Product  string
	Variant  string
}

func parseInitializationFlags(c *cli.Context) updater.State {
	return updater.State{
		Vendor:   c.String("vendor"),
		Product:  c.String("product"),
		Variant:  c.String("variant"),
		ClientId: c.String("client-id"),
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
	}
}

func createInitializationFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "vendor", Usage: "Vendor name", Required: true},
		&cli.StringFlag{Name: "product", Usage: "Product name", Required: true},
		&cli.StringFlag{Name: "variant", Usage: "Product variant", Required: true},
		&cli.StringFlag{Name: "client-id", Usage: "Optional: Client ID, if you already have one"},
	}
}
