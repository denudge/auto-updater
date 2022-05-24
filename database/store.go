package database

import (
	"context"
	"fmt"
	"github.com/denudge/auto-updater/catalog"
	"github.com/uptrace/bun"
	"strings"
	"time"
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
	dbApp, err := store.getApp(app.Vendor, app.Product, false)
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

	return store.getApp(app.Vendor, app.Product, true)
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
		if err != nil {
			return nil, fmt.Errorf("cannot find group \"%s\"", group)
		}
		groupMap[group.Name] = group
	}

	return groupMap, nil
}

func (store *DbCatalogStore) FindApp(vendor string, product string) (*catalog.App, error) {
	dbApp, err := store.getApp(vendor, product, true)
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

func (store *DbCatalogStore) LatestReleases(limit int) ([]*catalog.Release, error) {
	releases := make([]Release, 0, limit)

	query := store.db.NewSelect().
		Model(&releases).
		Relation("App").
		Relation("Groups").
		OrderExpr("id DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Scan(store.ctx)

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

	dbApp, err := store.getApp(group.App.Vendor, group.App.Product, false)
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

	dbApp, err := store.getApp(release.App.Vendor, release.App.Product, true)
	if err != nil {
		return nil, err
	}

	r.AppId = dbApp.Id

	// Check groups before inserting any release
	groups := make([]Group, 0)
	if release.Groups == nil || len(release.Groups) < 1 {
		// Connect default groups
		groups = dbApp.GetDefaultGroups()
	} else {
		groups, err = store.getGroups(dbApp, release.Groups)
		if err != nil {
			return nil, err
		}
	}

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

	linkGroups := false
	if _, err := stmt.Exec(store.ctx); err != nil {
		// Instead of returning an error, we just return the older release
		// So the caller can determine the release was already there
		if !strings.Contains(err.Error(), "violates unique constraint") {
			return nil, err
		}
	} else if dbApp.Groups != nil && len(dbApp.Groups) > 0 {
		linkGroups = true
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

	if linkGroups {
		for _, group := range groups {
			groupRelation := ReleaseToGroup{
				GroupId:   group.Id,
				ReleaseId: stored.Id,
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
			Where("release.id = ?", stored.Id).
			Relation("App").
			Relation("Groups").
			Scan(store.ctx)

		if err != nil {
			return nil, err
		}
	}

	return stored.ToCatalogRelease(), nil
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

func (store *DbCatalogStore) getApp(vendor string, product string, withGroups bool) (*App, error) {
	app := App{}

	query := store.db.NewSelect().
		Model(&app).
		Where("vendor = ?", vendor).
		Where("product = ?", product)

	if withGroups {
		query.Relation("Groups")
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

func (store *DbCatalogStore) ListGroups(filter catalog.GroupFilter, limit int) ([]*catalog.Group, error) {
	// Reserve reasonable space
	groups := make([]Group, 0, limit)

	dbApp, err := store.getApp(filter.Vendor, filter.Product, false)
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

func (store *DbCatalogStore) FetchReleases(filter catalog.Filter, limit int) ([]*catalog.Release, error) {
	dbApp, err := store.getApp(filter.Vendor, filter.Product, false)
	if err != nil || dbApp == nil {
		return []*catalog.Release{}, err
	}

	// Reserve at least some reasonable space
	releases := make([]Release, 0, 16)

	query := store.db.NewSelect().
		Model(&releases).
		Relation("App").
		Relation("Groups").
		Where("release.app_id = ?", dbApp.Id)

	query, err = store.filterQuery(dbApp, query, filter)
	if err != nil {
		return []*catalog.Release{}, err
	}

	query.OrderExpr("id DESC")

	if limit > 0 {
		query.Limit(limit)
	}

	err = query.
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

	return store.FetchReleases(filter, 0)
}

func (store *DbCatalogStore) SetStability(filter catalog.Filter, stability bool) ([]*catalog.Release, error) {
	// TODO: Implement update functionality

	return store.FetchReleases(filter, 0)
}

func (store *DbCatalogStore) SetUpgradeTarget(filter catalog.Filter, target catalog.UpgradeTarget) ([]*catalog.Release, error) {
	// TODO: Implement update functionality

	return store.FetchReleases(filter, 0)
}

func (store *DbCatalogStore) filterQuery(app *App, stmt *bun.SelectQuery, filter catalog.Filter) (*bun.SelectQuery, error) {
	if filter.Variant != "" {
		stmt.Where("variant = ?", filter.Variant)
	}

	if filter.OS != "" {
		stmt.Where("os = ?", filter.OS)
	}

	if filter.Arch != "" {
		stmt.Where("arch = ?", filter.Arch)
	}

	if filter.Alias != "" {
		stmt.Where("alias = ?", filter.Alias)
	}

	if !filter.WithUnstable {
		stmt.Where("unstable = false")
	}

	if filter.Groups != nil && len(filter.Groups) > 0 {
		stmt.Join("LEFT JOIN releases_groups AS rg ON rg.release_id = release.id")
		if len(filter.Groups) == 1 && filter.Groups[0] == "public" {
			stmt.Where("rg.group_id IS NULL")
		} else {
			groups, err := store.getGroups(app, filter.Groups)
			if err != nil {
				return nil, err
			}

			groupIds := make([]int64, len(groups))
			for i, groupObj := range groups {
				groupIds[i] = groupObj.Id
			}

			stmt.Where("(rg.group_id IN (?) OR rg.group_id IS NULL)", bun.In(groupIds))
		}
	}

	return stmt, nil
}
