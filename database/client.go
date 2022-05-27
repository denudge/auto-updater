package database

import (
	"github.com/denudge/auto-updater/catalog"
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
