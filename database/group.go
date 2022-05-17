package database

import (
	"github.com/denudge/auto-updater/catalog"
	"github.com/uptrace/bun"
	"time"
)

type Group struct {
	bun.BaseModel `bun:"table:groups"`
	Id            int64     `bun:"id,pk,autoincrement"`
	AppId         int64     `bun:"app_id"`
	App           *App      `bun:"rel:belongs-to,join:app_id=id"`
	Name          string    `bun:"name"`
	CreatedAt     time.Time `bun:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at"`
}

func (g *Group) ToCatalogGroup() *catalog.Group {
	return &catalog.Group{
		App: &catalog.App{
			Vendor:  g.App.Vendor,
			Product: g.App.Product,
			Name:    g.App.Name,
		},
		Name:    g.Name,
		Created: g.CreatedAt,
		Updated: g.UpdatedAt,
	}
}

func FromCatalogGroup(g *catalog.Group) Group {
	app := FromCatalogApp(g.App)

	return Group{
		App:       &app,
		Name:      g.Name,
		CreatedAt: g.Created,
		UpdatedAt: time.Now(),
	}
}

func transformGroups(groups []Group) ([]*catalog.Group, error) {
	if len(groups) < 1 {
		return []*catalog.Group{}, nil
	}

	out := make([]*catalog.Group, len(groups))

	for i, group := range groups {
		out[i] = group.ToCatalogGroup()
	}

	return out, nil
}
