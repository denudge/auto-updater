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

	registerRelations(db)

	return db
}

func Close(db *bun.DB) {
	err := db.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func registerRelations(db *bun.DB) {
	// Register many to many model so bun can better recognize m2m relation.
	// This should be done before you use the models for the first time.
	db.RegisterModel((*ReleaseToGroup)(nil))
	db.RegisterModel((*ClientToGroup)(nil))
	db.RegisterModel((*VariantToDefaultGroup)(nil))
}
