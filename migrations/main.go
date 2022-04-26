package migrations

import "github.com/uptrace/bun/migrate"

var Migrations = migrate.NewMigrations()

// Source: https://github.com/uptrace/bun/blob/78c14c8b92e1c5091a0d06f49835fabc95b879a8/example/migrate/migrations/main.go
func init() {
	if err := Migrations.DiscoverCaller(); err != nil {
		panic(err)
	}
}
