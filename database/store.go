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

func (store *DbCatalogStore) LatestReleases(limit int) ([]*catalog.Release, error) {
	releases := make([]Release, 0, limit)

	err := store.db.NewSelect().
		Model(&releases).
		OrderExpr("id DESC").
		Limit(limit).
		Scan(store.ctx)

	if err != nil {
		return []*catalog.Release{}, err
	}

	return transformReleases(releases)
}

func (store *DbCatalogStore) Store(release *catalog.Release, allowUpdate bool) (*catalog.Release, error) {
	r := FromCatalogRelease(release)

	stmt := store.db.NewInsert().
		Model(&r)

	if allowUpdate {
		stmt = stmt.
			// All fields except the ones in the unique constraint (as well as ID and date)
			On("CONFLICT ON CONSTRAINT releases_ux DO UPDATE").
			Set("name = EXCLUDED.name").
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
	err := store.db.NewSelect().
		Model(&stored).
		Where("vendor = ?", release.Vendor).
		Where("product = ?", release.Product).
		Where("variant = ?", release.Variant).
		Where("os = ?", release.OS).
		Where("arch = ?", release.Arch).
		Where("version = ?", release.Version).
		Scan(store.ctx)

	if err != nil {
		return nil, err
	}

	return stored.ToCatalogRelease(), nil
}

func (store *DbCatalogStore) Fetch(filter catalog.Filter) ([]*catalog.Release, error) {
	// Reserve at least some reasonable space
	releases := make([]Release, 0, 16)

	err := filterQuery(
		store.db.NewSelect().
			Model(&releases).
			Where("vendor = ?", filter.Vendor).
			Where("product = ?", filter.Product),
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

	return store.Fetch(filter)
}

func (store *DbCatalogStore) SetStability(filter catalog.Filter, stability bool) ([]*catalog.Release, error) {
	// TODO: Implement update functionality

	return store.Fetch(filter)
}

func (store *DbCatalogStore) SetUpgradeTarget(filter catalog.Filter, target catalog.UpgradeTarget) ([]*catalog.Release, error) {
	// TODO: Implement update functionality

	return store.Fetch(filter)
}

func transformReleases(releases []Release) ([]*catalog.Release, error) {
	if len(releases) < 1 {
		return []*catalog.Release{}, nil
	}

	out := make([]*catalog.Release, len(releases))

	for i, release := range releases {
		out[i] = release.ToCatalogRelease()
	}

	return out, nil
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
