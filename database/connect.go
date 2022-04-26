package database

import (
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"log"
)

func Connect(dsn string) *bun.DB {
	pgsql := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	return bun.NewDB(pgsql, pgdialect.New())
}

func Close(db *bun.DB) {
	err := db.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
}
