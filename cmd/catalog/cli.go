package main

import (
	"fmt"
	"github.com/denudge/auto-updater/database"
	"github.com/denudge/auto-updater/migrations"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
)

func NewCli(app *App, api *Api) *cli.App {
	return &cli.App{
		Name: "catalog",
		Commands: []*cli.Command{
			// A bunch of database (migration) related commands
			database.NewCommand(migrate.NewMigrator(app.db, migrations.Migrations)),
			{
				Name:  "serve",
				Usage: "run HTTP catalog API server",
				Action: func(c *cli.Context) error {
					api.Serve()
					return nil
				},
			},
			app.createAppCommands(),
			app.createVariantCommands(),
			app.createGroupCommands(),
			app.createReleaseCommands(),
			app.createClientCommands(),
		},
	}
}

func createLimitFlag(defaultValue int) []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:        "limit",
			Usage:       fmt.Sprintf("Optional: Limit result set. Default: %d", defaultValue),
			DefaultText: fmt.Sprintf("%d", defaultValue),
		},
	}
}

func parseLimitFlag(c *cli.Context, defaultValue int) int {
	given := c.String("limit")
	if given == "" {
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
