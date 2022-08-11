package database

import (
	"fmt"
	"github.com/denudge/auto-updater/catalog"
	"github.com/uptrace/bun"
	"strings"
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
		App:     g.App.ToCatalogApp(),
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

func (store *DbCatalogStore) ListGroups(filter catalog.GroupFilter, limit int) ([]*catalog.Group, error) {
	// Reserve reasonable space
	groups := make([]Group, 0, limit)

	dbApp, err := store.getApp(filter.Vendor, filter.Product, false, false)
	if err != nil || dbApp == nil {
		return []*catalog.Group{}, err
	}

	query := store.db.NewSelect().
		Model(&groups).
		Relation("App").
		Where("app_id = ?", dbApp.Id)

	if filter.Name != "" {
		query = query.Where("name = ?", filter.Name)
	}

	query = query.
		OrderExpr("\"name\" ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err = query.Scan(store.ctx)

	if err != nil {
		return []*catalog.Group{}, err
	}

	return transformGroups(groups)
}

func (store *DbCatalogStore) StoreGroup(
	group *catalog.Group,
	allowUpdate bool,
) (*catalog.Group, error) {
	g := FromCatalogGroup(group)

	dbApp, err := store.getApp(group.App.Vendor, group.App.Product, false, false)
	if err != nil {
		return nil, err
	}

	g.AppId = dbApp.Id

	stmt := store.db.NewInsert().
		Model(&g)

	if allowUpdate {
		stmt = stmt.
			// All fields except the ones in the unique constraint (as well as ID and date)
			On("CONFLICT ON CONSTRAINT groups_ux IGNORE")
	}

	if _, err := stmt.Exec(store.ctx); err != nil {
		// Instead of returning an error, we just return the older release
		// So the caller can determine the group was already there
		if !strings.Contains(err.Error(), "violates unique constraint") {
			return nil, err
		}
	}

	stored := Group{}

	err = store.db.NewSelect().
		Model(&stored).
		Where("app_id = ?", dbApp.Id).
		Where("\"group\".\"name\" = ?", group.Name).
		Relation("App").
		Scan(store.ctx)

	if err != nil {
		return nil, err
	}

	return stored.ToCatalogGroup(), nil
}

func (store *DbCatalogStore) getGroups(app *App, groupNames []string) ([]Group, error) {
	if groupNames == nil || len(groupNames) < 1 {
		return []Group{}, nil
	}

	groups := make([]Group, 0, len(groupNames))
	groupMap, err := store.getGroupMap(app.Id)
	if err != nil {
		return nil, err
	}

	for _, group := range groupNames {
		if group == "public" {
			continue
		}

		groupObj, ok := groupMap[group]
		if !ok {
			return nil, fmt.Errorf("unknown group \"%s\"", group)
		}

		groups = append(groups, groupObj)
	}

	return groups, nil
}

func (store *DbCatalogStore) getGroupMap(appId int64) (map[string]Group, error) {
	groups := make([]Group, 0, 4)

	err := store.db.NewSelect().
		Model(&groups).
		Where("app_id = ?", appId).
		OrderExpr("id DESC").
		Scan(store.ctx)

	if err != nil {
		return nil, err
	}

	groupMap := make(map[string]Group)
	for _, group := range groups {
		groupMap[group.Name] = group
	}

	return groupMap, nil
}

func (store *DbCatalogStore) getGroup(appId int64, name string) (*Group, error) {
	group := Group{}

	fmt.Printf("Searching for group %s\n", name)

	err := store.db.NewSelect().
		Model(&group).
		Where("app_id = ?", appId).
		Where("name = ?", name).
		Scan(store.ctx)

	if err != nil {
		return nil, err
	}

	return &group, nil
}
