package database

import (
	"github.com/denudge/auto-updater/catalog"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Client struct {
	bun.BaseModel `bun:"table:clients"`
	Id            int64     `bun:"id,pk,autoincrement"`
	AppId         int64     `bun:"app_id"`
	App           *App      `bun:"rel:belongs-to,join:app_id=id"`
	Variant       string    `bun:"variant"` // will be enforced when fetching releases
	Uuid          string    `bun:"uuid"`
	Name          string    `bun:"name"`
	Active        bool      `bun:"active,default:true"`
	Locked        bool      `bun:"locked,default:false"`
	CreatedAt     time.Time `bun:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at"`
	Groups        []Group   `bun:"m2m:clients_groups,join:Client=Group"`
}

func (c *Client) ToCatalogClient() *catalog.Client {
	client := &catalog.Client{
		App: &catalog.App{
			Vendor:  c.App.Vendor,
			Product: c.App.Product,
			Name:    c.App.Name,
		},
		Variant: c.Variant,
		Uuid:    c.Uuid,
		Name:    c.Name,
		Active:  c.Active,
		Locked:  c.Locked,
		Created: c.CreatedAt,
		Updated: c.UpdatedAt,
	}

	// Transform groups to simple strings
	client.Groups = make([]string, len(c.Groups))
	for i, group := range c.Groups {
		client.Groups[i] = group.Name
	}

	return client
}

func (store *DbCatalogStore) RegisterClient(app *catalog.App, variant string, groups []string) (*catalog.Client, error) {
	dbApp, err := store.getApp(app.Vendor, app.Product, false)
	if err != nil {
		return nil, err
	}

	// Make sure we have the right default groups
	groupNames := make([]string, 0)
	groupObjs := make([]Group, 0)
	if groups != nil && len(groups) > 0 {
		groupNames = groups

		groupObjs, err = store.getGroups(dbApp, groupNames)
		if err != nil {
			return nil, err
		}
	}

	client := Client{
		AppId:     dbApp.Id,
		Variant:   variant,
		Uuid:      uuid.NewString(),
		Active:    true,
		Locked:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	stmt := store.db.NewInsert().
		Model(&client)

	if _, err := stmt.Exec(store.ctx); err != nil {
		return nil, err
	}

	stored := Client{}

	err = store.db.NewSelect().
		Model(&stored).
		Where("app_id = ?", dbApp.Id).
		Where("uuid = ?", client.Uuid).
		Relation("App").
		Relation("Groups").
		Scan(store.ctx)

	if err != nil {
		return nil, err
	}

	if len(groupObjs) > 0 {
		for _, group := range groupObjs {
			groupRelation := ClientToGroup{
				GroupId:   group.Id,
				ClientId:  stored.Id,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if _, err = store.db.NewInsert().
				Model(&groupRelation).
				Exec(store.ctx); err != nil {
				return nil, err
			}
		}

		err = store.db.NewSelect().
			Model(&stored).
			Where("app_id = ?", dbApp.Id).
			Where("uuid = ?", client.Uuid).
			Relation("App").
			Relation("Groups").
			Scan(store.ctx)

		if err != nil {
			return nil, err
		}
	}

	return stored.ToCatalogClient(), nil
}
