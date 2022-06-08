package database

import (
	"github.com/denudge/auto-updater/catalog"
	"github.com/uptrace/bun"
	"strings"
	"time"
)

type Variant struct {
	bun.BaseModel `bun:"table:variants"`
	Id            int64     `bun:"id,pk,autoincrement"`
	AppId         int64     `bun:"app_id"`
	App           *App      `bun:"rel:belongs-to,join:app_id=id"`
	Name          string    `bun:"name"`
	Active        bool      `bun:"active,default:true"`
	Locked        bool      `bun:"locked,default:false"`
	AllowRegister bool      `bun:"allow_register,default:true"`
	UpgradeTarget string    `bun:"upgrade_target"` // If empty, the default upgrade target will be used
	CreatedAt     time.Time `bun:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at"`
	Groups        []Group   `bun:"m2m:variants_default_groups,join:Variant=Group"`
}

func (v *Variant) ToCatalogVariant() *catalog.Variant {
	return &catalog.Variant{
		App:           v.App.ToCatalogApp(),
		Name:          v.Name,
		Active:        v.Active,
		Locked:        v.Locked,
		AllowRegister: v.AllowRegister,
		UpgradeTarget: catalog.UpgradeTarget(v.UpgradeTarget),
		Created:       v.CreatedAt,
		Updated:       v.UpdatedAt,
	}
}

func FromCatalogVariant(v *catalog.Variant) Variant {
	app := FromCatalogApp(v.App)

	return Variant{
		App:           &app,
		Name:          v.Name,
		Active:        app.Active,
		Locked:        app.Locked,
		AllowRegister: app.AllowRegister,
		UpgradeTarget: string(app.UpgradeTarget),
		CreatedAt:     v.Created,
		UpdatedAt:     time.Now(),
	}
}

func transformVariants(variants []Variant) ([]*catalog.Variant, error) {
	if len(variants) < 1 {
		return []*catalog.Variant{}, nil
	}

	out := make([]*catalog.Variant, len(variants))

	for i, variant := range variants {
		out[i] = variant.ToCatalogVariant()
	}

	return out, nil
}

type VariantToDefaultGroup struct {
	bun.BaseModel `bun:"table:variants_default_groups"`
	Id            int64     `bun:"id,pk,autoincrement"`
	VariantId     int64     `bun:"variant_id"`
	Variant       *Variant  `bun:"rel:belongs-to,join:variant_id=id"`
	GroupId       int64     `bun:"group_id"`
	Group         *Group    `bun:"rel:belongs-to,join:group_id=id"`
	CreatedAt     time.Time `bun:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at"`
}

func (store *DbCatalogStore) ListVariants(filter catalog.VariantFilter, limit int) ([]*catalog.Variant, error) {
	// Reserve reasonable space
	variants := make([]Variant, 0, limit)

	dbApp, err := store.getApp(filter.Vendor, filter.Product, false, false)
	if err != nil || dbApp == nil {
		return []*catalog.Variant{}, err
	}

	query := store.db.NewSelect().
		Model(&variants).
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
		return []*catalog.Variant{}, err
	}

	return transformVariants(variants)
}

func (store *DbCatalogStore) StoreVariant(
	variant *catalog.Variant,
	allowUpdate bool,
) (*catalog.Variant, error) {
	v := FromCatalogVariant(variant)

	dbApp, err := store.getApp(variant.App.Vendor, variant.App.Product, false, false)
	if err != nil {
		return nil, err
	}

	v.AppId = dbApp.Id

	// Check default groups before inserting any variant
	defaultGroups := make([]Group, 0)
	if variant.DefaultGroups != nil && len(variant.DefaultGroups) > 0 {
		defaultGroups, err = store.getGroups(dbApp, variant.DefaultGroups)
		if err != nil {
			return nil, err
		}
	}

	stmt := store.db.NewInsert().
		Model(&v)

	linkGroups := false
	if allowUpdate {
		stmt = stmt.
			// All fields except the ones in the unique constraint (as well as ID and date)
			On("CONFLICT ON CONSTRAINT variants_ux IGNORE")
	} else if len(defaultGroups) > 0 {
		linkGroups = true
	}

	if _, err := stmt.Exec(store.ctx); err != nil {
		// Instead of returning an error, we just return the older release
		// So the caller can determine the variant was already there
		if !strings.Contains(err.Error(), "violates unique constraint") {
			return nil, err
		}
	}

	stored := Variant{}

	err = store.db.NewSelect().
		Model(&stored).
		Where("app_id = ?", dbApp.Id).
		Where("\"variant\".\"name\" = ?", variant.Name).
		Relation("App").
		Scan(store.ctx)

	if err != nil {
		return nil, err
	}

	if linkGroups {
		for _, group := range defaultGroups {
			groupRelation := VariantToDefaultGroup{
				GroupId:   group.Id,
				VariantId: stored.Id,
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
			Where("\"variant\".\"name\" = ?", variant.Name).
			Relation("App").
			Scan(store.ctx)

		if err != nil {
			return nil, err
		}
	}

	return stored.ToCatalogVariant(), nil
}
