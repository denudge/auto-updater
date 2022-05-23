package database

import (
	"github.com/denudge/auto-updater/catalog"
	"github.com/uptrace/bun"
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
	UpgradeTarget string    `bun:"upgrade_target"` // If empty, the default upgrade target will be used
	CreatedAt     time.Time `bun:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at"`
	Groups        []Group   `bun:"rel:has-many,join:id=app_id"`
}

func (app *App) ToCatalogApp() *catalog.App {
	a := &catalog.App{
		Vendor:        app.Vendor,
		Product:       app.Product,
		Name:          app.Name,
		Active:        app.Active,
		Locked:        app.Locked,
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

	return a
}

func FromCatalogApp(app *catalog.App) App {
	return App{
		Vendor:        app.Vendor,
		Product:       app.Product,
		Name:          app.Name,
		Active:        app.Active,
		Locked:        app.Locked,
		UpgradeTarget: string(app.UpgradeTarget),
		CreatedAt:     app.Created,
		UpdatedAt:     app.Updated,
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
