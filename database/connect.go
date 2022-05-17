package database

import (
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"log"
)

func Connect(dsn string) *bun.DB {
	pgsql := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(pgsql, pgdialect.New())

	debug := bundebug.NewQueryHook(
		// disable the hook
		bundebug.WithEnabled(false),

		bundebug.WithVerbose(true),

		// BUNDEBUG=1 logs failed queries
		// BUNDEBUG=2 logs all queries
		bundebug.FromEnv("BUNDEBUG"),
	)

	db.AddQueryHook(debug)

	return db
}

func Close(db *bun.DB) {
	err := db.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
}
