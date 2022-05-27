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
	Default       bool      `bun:"default"`
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

type ReleaseToGroup struct {
	bun.BaseModel `bun:"table:releases_groups"`
	Id            int64     `bun:"id,pk,autoincrement"`
	ReleaseId     int64     `bun:"release_id"`
	Release       *Release  `bun:"rel:belongs-to,join:release_id=id"`
	GroupId       int64     `bun:"group_id"`
	Group         *Group    `bun:"rel:belongs-to,join:group_id=id"`
	CreatedAt     time.Time `bun:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at"`
}

type ClientToGroup struct {
	bun.BaseModel `bun:"table:clients_groups"`
	Id            int64     `bun:"id,pk,autoincrement"`
	ClientId      int64     `bun:"client_id"`
	Client        *Client   `bun:"rel:belongs-to,join:client_id=id"`
	GroupId       int64     `bun:"group_id"`
	Group         *Group    `bun:"rel:belongs-to,join:group_id=id"`
	CreatedAt     time.Time `bun:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at"`
}
