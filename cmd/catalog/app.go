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
	store catalog.StoreInterface
}

func NewApp(db *bun.DB, ctx context.Context) *App {
	return &App{
		db:    db,
		store: database.NewDbCatalogStore(db, ctx),
	}
}

func (app *App) ListApps(limit int) error {
	dbStore, ok := app.store.(*database.DbCatalogStore)
	if !ok {
		fmt.Printf("Cannot print apps")
		return nil
	}

	latest, err := dbStore.ListApps(limit)
	if err != nil {
		return err
	}

	for _, release := range latest {
		fmt.Printf("%s\n", release)
	}

	return nil
}

func (app *App) ListLatestReleases(limit int) error {
	dbStore, ok := app.store.(*database.DbCatalogStore)
	if !ok {
		fmt.Printf("Cannot print latest releases")
		return nil
	}

	latest, err := dbStore.LatestReleases(limit)
	if err != nil {
		return err
	}

	for _, release := range latest {
		fmt.Printf("%s\n", release)
	}

	return nil
}
