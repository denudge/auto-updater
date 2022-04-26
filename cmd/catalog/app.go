package main

import (
	"context"
	"fmt"
	"github.com/denudge/auto-updater/catalog"
	"github.com/denudge/auto-updater/database"
	"github.com/uptrace/bun"
)

type App struct {
	db    *bun.DB
	store *database.DbCatalogStore
}

func NewApp(db *bun.DB, ctx context.Context) *App {
	return &App{
		db:    db,
		store: database.NewDbCatalogStore(db, context.Background()),
	}
}

func (app *App) ListLatestReleases() error {
	latest, err := app.store.LatestReleases(10)
	if err != nil {
		return err
	}

	for _, release := range latest {
		fmt.Printf("%s\n", release)
	}

	return nil
}

func (app *App) Fetch(filter catalog.Filter) ([]*catalog.Release, error) {
	return app.store.Fetch(filter)
}
