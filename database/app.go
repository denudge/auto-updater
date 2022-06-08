package database

import (
	"fmt"
	"github.com/denudge/auto-updater/catalog"
	"github.com/uptrace/bun"
	"strings"
	"time"
)

type App struct {
	bun.BaseModel `bun:"table:apps"`
	Id            int64     `bun:"id,pk,autoincrement"`
	Vendor        string    `bun:"vendor"`
	Product       string    `bun:"product"`
	Name          string    `bun:"name"`
	Active        bool      `bun:"active,default:true"`
	Locked        bool      `bun:"locked,default:false"`
	AllowRegister bool      `bun:"allow_register,default:true"`
	UpgradeTarget string    `bun:"upgrade_target"` // If empty, the default upgrade target will be used
	CreatedAt     time.Time `bun:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at"`
	Groups        []Group   `bun:"rel:has-many,join:id=app_id"`
	Variants      []Variant `bun:"rel:has-many,join:id=app_id"`
}

func (app *App) ToCatalogApp() *catalog.App {
	a := &catalog.App{
		Vendor:        app.Vendor,
		Product:       app.Product,
		Name:          app.Name,
		Active:        app.Active,
		Locked:        app.Locked,
		AllowRegister: app.AllowRegister,
		UpgradeTarget: catalog.UpgradeTarget(app.UpgradeTarget),
		Created:       app.CreatedAt,
		Updated:       app.UpdatedAt,
	}

	if app.Groups != nil {
		a.Groups = make([]string, len(app.Groups))
		a.DefaultGroups = make([]string, 0, len(app.Groups))
		for i, group := range app.Groups {
			a.Groups[i] = group.Name
			if group.Default {
				a.DefaultGroups = append(a.DefaultGroups, group.Name)
			}
		}
	}

	if app.Variants != nil {
		a.Variants = make([]string, len(app.Variants))
		for i, variant := range app.Variants {
			a.Variants[i] = variant.Name
		}
	}

	return a
}

func (app *App) GetDefaultGroups() []Group {
	if app.Groups == nil || len(app.Groups) < 1 {
		return []Group{}
	}

	defaultGroups := make([]Group, 0, len(app.Groups))
	for _, group := range app.Groups {
		if group.Default {
			defaultGroups = append(defaultGroups, group)
		}
	}

	return defaultGroups
}

func FromCatalogApp(app *catalog.App) App {
	return App{
		Vendor:        app.Vendor,
		Product:       app.Product,
		Name:          app.Name,
		Active:        app.Active,
		Locked:        app.Locked,
		AllowRegister: app.AllowRegister,
		UpgradeTarget: string(app.UpgradeTarget),
		CreatedAt:     app.Created,
		UpdatedAt:     time.Now(),
	}
}

func transformApps(apps []App) ([]*catalog.App, error) {
	if len(apps) < 1 {
		return []*catalog.App{}, nil
	}

	out := make([]*catalog.App, len(apps))

	for i, app := range apps {
		out[i] = app.ToCatalogApp()
	}

	return out, nil
}

func (store *DbCatalogStore) FindApp(vendor string, product string) (*catalog.App, error) {
	dbApp, err := store.getApp(vendor, product, true, true)
	if err != nil {
		return nil, err
	}

	return dbApp.ToCatalogApp(), nil
}

func (store *DbCatalogStore) ListApps(limit int) ([]*catalog.App, error) {
	apps := make([]App, 0, limit)

	query := store.db.NewSelect().
		Model(&apps).
		Relation("Groups").
		Relation("Variants").
		OrderExpr("id DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Scan(store.ctx)

	if err != nil {
		return []*catalog.App{}, err
	}

	return transformApps(apps)
}

func (store *DbCatalogStore) StoreApp(app *catalog.App, allowUpdate bool) (*catalog.App, error) {
	a := FromCatalogApp(app)

	stmt := store.db.NewInsert().
		Model(&a)

	if allowUpdate {
		stmt = stmt.
			// All fields except the ones in the unique constraint (as well as ID and date)
			On("CONFLICT ON CONSTRAINT apps_ux DO UPDATE").
			Set("name = EXCLUDED.name")
	}

	if _, err := stmt.Exec(store.ctx); err != nil {
		// Instead of returning an error, we just return the older app
		// So the caller can determine the app was already there
		if !strings.Contains(err.Error(), "violates unique constraint") {
			return nil, err
		}
	}

	return store.FindApp(app.Vendor, app.Product)
}

func (store *DbCatalogStore) SetAppDefaultGroups(app *catalog.App) (*catalog.App, error) {
	// Now make sure we have the right default groups
	dbApp, err := store.getApp(app.Vendor, app.Product, false, false)
	if err != nil {
		return nil, err
	}

	groupNames := make([]string, 0)
	if app.DefaultGroups != nil && len(app.DefaultGroups) > 0 {
		groupNames = app.DefaultGroups
	}

	dbApp, err = store.setAppDefaultGroups(dbApp, groupNames)
	if err != nil {
		return nil, err
	}

	return dbApp.ToCatalogApp(), nil
}

func (store *DbCatalogStore) setAppDefaultGroups(app *App, groups []string) (*App, error) {
	_, err := store.getGroups(app, groups)
	if err != nil {
		return nil, err
	}

	// Deactivate all previous default groups
	_, err = store.db.NewUpdate().
		Model(&Group{}).
		Where("app_id = ?", app.Id).
		Set("\"default\" = false").
		Exec(store.ctx)

	if err != nil {
		return nil, err
	}

	// Set new default groups
	_, err = store.db.NewUpdate().
		Model(&Group{}).
		Where("app_id = ?", app.Id).
		Where("\"name\" IN (?)", bun.In(groups)).
		Set("\"default\" = true").
		Exec(store.ctx)

	if err != nil {
		return nil, err
	}

	return store.getApp(app.Vendor, app.Product, true, true)
}

func (store *DbCatalogStore) getApp(vendor string, product string, withGroups bool, withVariants bool) (*App, error) {
	app := App{}

	query := store.db.NewSelect().
		Model(&app).
		Where("vendor = ?", vendor).
		Where("product = ?", product)

	if withGroups {
		query.Relation("Groups")
	}

	if withVariants {
		query.Relation("Variants")
	}

	err := query.
		Scan(store.ctx)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("unknown app \"%s %s\"", vendor, product)
		}

		return nil, err
	}

	return &app, nil
}
