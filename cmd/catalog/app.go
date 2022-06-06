package main

import (
	"context"
	"fmt"
	"github.com/denudge/auto-updater/catalog"
	"github.com/denudge/auto-updater/database"
	"github.com/uptrace/bun"
)

// App implements the Catalog interface and provides the user-facing server part as well (as some basic management methods)
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

func (app *App) RegisterClient(vendor string, product string, variant string) (*catalog.ClientState, error) {
	dbApp, err := app.store.FindApp(vendor, product)
	if err != nil {
		return nil, err
	}

	client, err := app.store.RegisterClient(dbApp, variant, []string{})
	if err != nil {
		return nil, err
	}

	state := &catalog.ClientState{
		ClientId: client.Uuid,
		Vendor:   vendor,
		Product:  product,
		Variant:  variant,
	}

	return state, nil
}

func (app *App) ShouldUpgrade(state *catalog.ClientState) (catalog.Criticality, error) {
	if !state.IsValid() {
		return catalog.None, fmt.Errorf("state is not valid. Please register first")
	}

	if !state.IsInstalled() {
		return catalog.None, nil
	}

	filter := catalog.Filter{
		Vendor:         state.Vendor,
		Product:        state.Product,
		Variant:        state.Variant,
		EnforceVariant: true,
		MinVersion:     state.CurrentVersion,
	}

	releases, err := app.store.FetchReleases(filter, 0)
	if err != nil {
		return catalog.None, err
	}

	step, err := catalog.FindNextUpgrade(releases, state.CurrentVersion)
	if err != nil {
		return catalog.None, err
	}

	if step == nil {
		return catalog.None, nil
	}

	return step.Criticality, nil
}

func (app *App) FindNextUpgrade(state *catalog.ClientState) (*catalog.UpgradeStep, error) {
	if !state.IsValid() {
		return nil, fmt.Errorf("state is not valid. Please register first")
	}

	filter := state.ToFilter()

	releases, err := app.store.FetchReleases(filter, 0)
	if err != nil {
		return nil, err
	}

	// If no version is installed yet, use the "install" version
	if !state.IsInstalled() {
		return catalog.FindInstallVersion(releases, state.WithUnstable)
	}

	return catalog.FindNextUpgrade(releases, state.CurrentVersion)
}

func (app *App) FindUpgradePath(state *catalog.ClientState) (*catalog.UpgradePath, error) {
	if !state.IsValid() {
		return nil, fmt.Errorf("state is not valid. Please register first")
	}

	filter := state.ToFilter()

	releases, err := app.store.FetchReleases(filter, 0)
	if err != nil {
		return nil, err
	}

	// If no version is installed yet, use the "install" version
	if !state.IsInstalled() {
		step, err := catalog.FindInstallVersion(releases, state.WithUnstable)
		if err != nil {
			return nil, err
		}

		return step.ToPath(), err
	}

	return catalog.FindUpgradePath(releases, state.CurrentVersion)
}

// ListLatestReleases is an internal function for server management
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
