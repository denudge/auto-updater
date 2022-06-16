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

func newCli(updaterApp *app.Updater) *cli.App {
	return &cli.App{
		Name:  "updater",
		Usage: "the auto-updater client console",
		Flags: createGlobalFlags(),
		Before: func(c *cli.Context) error {
			updaterApp.Init(c.String("config-file"), c.String("server-address"))

			return updaterApp.CheckServerConfiguration()
		},
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "initialize the updater client",
				Flags: createInitializationFlags(),
				Action: func(c *cli.Context) error {
					if updaterApp.State != nil {
						return fmt.Errorf("config file already exists")
					}

					state := parseInitializationFlags(c)
					state.Server = updaterApp.BaseUrl

					if state.ClientId == "" {
						// We need to register
						response, err := updaterApp.Client.RegisterClient(state.Vendor, state.Product, state.Variant)
						if err != nil {
							return err
						}

						if response == nil || response.ClientId == "" {
							return fmt.Errorf("could not register client (unknown reason)")
						}

						state.ClientId = response.ClientId

						fmt.Printf("Successfully registered with client ID %s\n", state.ClientId)
					}

					return updaterApp.SaveState(&state)
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

func createGlobalFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "server-address", Usage: "catalog server address", Required: false},
		&cli.StringFlag{Name: "config-file", Usage: "config file", Required: false},
	}
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
