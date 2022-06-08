package cli

import (
	"fmt"
	"github.com/denudge/auto-updater/cmd/catalog/api"
	"github.com/denudge/auto-updater/cmd/catalog/app"
	"github.com/denudge/auto-updater/database"
	"github.com/denudge/auto-updater/migrations"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
)

type Console struct {
	Cli *cli.App
	app *app.App
}

func NewConsole(app *app.App, api *api.Api) *Console {
	console := &Console{
		app: app,
	}

	console.Cli = &cli.App{
		Name: "catalog",
		Commands: []*cli.Command{
			// A bunch of database (migration) related commands
			database.NewCommand(migrate.NewMigrator(app.Db, migrations.Migrations)),
			{
				Name:  "serve",
				Usage: "run HTTP catalog API server",
				Action: func(c *cli.Context) error {
					api.Serve()
					return nil
				},
			},
			console.createAppCommands(),
			console.createVariantCommands(),
			console.createGroupCommands(),
			console.createReleaseCommands(),
			console.createClientCommands(),
		},
	}

	return console
}

func createLimitFlag(defaultValue int) []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:        "limit",
			Usage:       "Optional: Limit result set",
			DefaultText: fmt.Sprintf("%d", defaultValue),
		},
	}
}

func parseLimitFlag(c *cli.Context, defaultValue int) int {
	given := c.Int("limit")
	if given == 0 {
		return defaultValue
	}

	return c.Int("limit")
}

func checkArguments(c *cli.Context, command string, fields []string) error {
	for _, field := range fields {
		if c.String(field) == "" {
			_ = cli.ShowCommandHelp(c, command)
			requiredFields := formatRequiredFields(fields)
			return fmt.Errorf("Missing argument \"%s\". At least %s must be specified.\n", field, requiredFields)
		}
	}

	return nil
}

func formatRequiredFields(fields []string) string {
	if len(fields) < 1 {
		return ""
	}

	requiredFields := fields[0]

	for i := 1; i < len(fields); i++ {
		if i < len(fields)-1 {
			requiredFields += ", " + fields[i]
		} else {
			requiredFields += " and " + fields[i]
		}
	}

	return requiredFields
}
