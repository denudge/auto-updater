package database

import (
	"context"
	"github.com/denudge/auto-updater/catalog"
	"github.com/uptrace/bun"
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

func (store *DbCatalogStore) Store(release *catalog.Release) error {
	// TODO: Implement release storage

	return nil
}

func (store *DbCatalogStore) Fetch(filter catalog.Filter) ([]*catalog.Release, error) {
	releases := make([]Release, 0, 16)

	err := filterQuery(
		store.db.NewSelect().
			Model(&releases).
			Where("vendor = ?", filter.Vendor).
			Where("product = ?", filter.Product),
		filter,
	).
		OrderExpr("id ASC").
		Scan(store.ctx)

	if err != nil {
		return []*catalog.Release{}, err
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
