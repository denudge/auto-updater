package database

import (
	"context"
	"github.com/denudge/auto-updater/catalog"
	"github.com/uptrace/bun"
	"strings"
)

type DbCatalogStore struct {
	db  *bun.DB
	ctx context.Context
}

func NewDbCatalogStore(db *bun.DB, ctx context.Context) *DbCatalogStore {
	return &DbCatalogStore{
		db:  db,
		ctx: ctx,
	}
}

func (store *DbCatalogStore) CreateApp(app *catalog.App, allowUpdate bool) (*catalog.App, error) {
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

func (store *DbCatalogStore) FindApp(vendor string, product string) (*catalog.App, error) {
	dbApp, err := store.GetApp(vendor, product)
	if err != nil {
		return nil, err
	}

	return dbApp.ToCatalogApp(), nil
}

func (store *DbCatalogStore) ListApps(limit int) ([]*catalog.App, error) {
	apps := make([]App, 0, limit)

	err := store.db.NewSelect().
		Model(&apps).
		Relation("DefaultGroups").
		OrderExpr("id DESC").
		Limit(limit).
		Scan(store.ctx)

	if err != nil {
		return []*catalog.App{}, err
	}

	return transformApps(apps)
}

func (store *DbCatalogStore) LatestReleases(limit int) ([]*catalog.Release, error) {
	releases := make([]Release, 0, limit)

	err := store.db.NewSelect().
		Model(&releases).
		Relation("App").
		Relation("Groups").
		OrderExpr("id DESC").
		Limit(limit).
		Scan(store.ctx)

	if err != nil {
		return []*catalog.Release{}, err
	}

	return transformReleases(releases)
}

func (store *DbCatalogStore) StoreGroup(
	group *catalog.Group,
	allowUpdate bool,
) (*catalog.Group, error) {
	g := FromCatalogGroup(group)

	dbApp, err := store.GetApp(group.App.Vendor, group.App.Product)
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
		// So the caller can determine the release was already there
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

func (store *DbCatalogStore) StoreRelease(
	release *catalog.Release,
	allowUpdate bool,
) (*catalog.Release, error) {
	r := FromCatalogRelease(release)

	dbApp, err := store.GetApp(release.App.Vendor, release.App.Product)
	if err != nil {
		return nil, err
	}

	r.AppId = dbApp.Id

	stmt := store.db.NewInsert().
		Model(&r)

	if allowUpdate {
		stmt = stmt.
			// All fields except the ones in the unique constraint (as well as ID and date)
			On("CONFLICT ON CONSTRAINT releases_ux DO UPDATE").
			Set("description = EXCLUDED.description").
			Set("alias = EXCLUDED.alias").
			Set("unstable = EXCLUDED.unstable").
			Set("signature = EXCLUDED.signature").
			Set("tags = EXCLUDED.tags").
			Set("upgrade_target = EXCLUDED.upgrade_target").
			Set("should_upgrade = EXCLUDED.should_upgrade")
	}

	if _, err := stmt.Exec(store.ctx); err != nil {
		// Instead of returning an error, we just return the older release
		// So the caller can determine the release was already there
		if !strings.Contains(err.Error(), "violates unique constraint") {
			return nil, err
		}
	}

	stored := Release{}

	err = store.db.NewSelect().
		Model(&stored).
		Where("app_id = ?", dbApp.Id).
		Where("variant = ?", release.Variant).
		Where("os = ?", release.OS).
		Where("arch = ?", release.Arch).
		Where("version = ?", release.Version).
		Relation("App").
		Relation("Groups").
		Scan(store.ctx)

	if err != nil {
		return nil, err
	}

	return stored.ToCatalogRelease(), nil
}

func (store *DbCatalogStore) GetApp(vendor string, product string) (*App, error) {
	app := App{}

	err := store.db.NewSelect().
		Model(&app).
		Where("vendor = ?", vendor).
		Where("product = ?", product).
		Scan(store.ctx)

	if err != nil {
		return nil, err
	}

	return &app, nil
}

func (store *DbCatalogStore) ListGroups(filter catalog.GroupFilter, limit int) ([]*catalog.Group, error) {
	// Reserve reasonable space
	groups := make([]Group, 0, limit)

	dbApp, err := store.GetApp(filter.Vendor, filter.Product)
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

	err = query.
		OrderExpr("\"name\" ASC").
		Limit(limit).
		Scan(store.ctx)

	if err != nil {
		return []*catalog.Group{}, err
	}

	return transformGroups(groups)
}

func (store *DbCatalogStore) FetchReleases(filter catalog.Filter) ([]*catalog.Release, error) {
	// Reserve at least some reasonable space
	releases := make([]Release, 0, 16)

	dbApp, err := store.GetApp(filter.Vendor, filter.Product)
	if err != nil || dbApp == nil {
		return []*catalog.Release{}, err
	}

	err = filterQuery(
		store.db.NewSelect().
			Model(&releases).
			Relation("App").
			Relation("Groups").
			Where("app_id = ?", dbApp.Id),
		filter,
	).
		OrderExpr("id DESC").
		Scan(store.ctx)

	if err != nil {
		return []*catalog.Release{}, err
	}

	if filter.FiltersVersions() {
		filtered := make([]Release, 0, len(releases))

		for _, release := range releases {
			if !filter.MatchVersion(release.Version) {
				continue
			}

			filtered = append(filtered, release)
		}

		releases = filtered
	}

	return transformReleases(releases)
}

func (store *DbCatalogStore) SetCriticality(filter catalog.Filter, criticality catalog.Criticality) ([]*catalog.Release, error) {
	// TODO: Implement update functionality

	return store.FetchReleases(filter)
}

func (store *DbCatalogStore) SetStability(filter catalog.Filter, stability bool) ([]*catalog.Release, error) {
	// TODO: Implement update functionality

	return store.FetchReleases(filter)
}

func (store *DbCatalogStore) SetUpgradeTarget(filter catalog.Filter, target catalog.UpgradeTarget) ([]*catalog.Release, error) {
	// TODO: Implement update functionality

	return store.FetchReleases(filter)
}

func filterQuery(stmt *bun.SelectQuery, filter catalog.Filter) *bun.SelectQuery {
	if filter.Variant != "" {
		stmt = stmt.
			Where("variant = ?", filter.Variant)
	}

	if filter.OS != "" {
		stmt = stmt.
			Where("os = ?", filter.OS)
	}

	if filter.Arch != "" {
		stmt = stmt.
			Where("arch = ?", filter.Arch)
	}

	if filter.Alias != "" {
		stmt = stmt.
			Where("alias = ?", filter.Alias)
	}

	if !filter.WithUnstable {
		stmt = stmt.
			Where("unstable = false")
	}

	return stmt
}
