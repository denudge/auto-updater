package database

import (
	"fmt"
	"github.com/denudge/auto-updater/catalog"
	"github.com/uptrace/bun"
	"strings"
	"time"
)

type Release struct {
	bun.BaseModel `bun:"table:releases"`
	Id            int64     `bun:"id,pk,autoincrement"`
	AppId         int64     `bun:"app_id"`
	App           *App      `bun:"rel:belongs-to,join:app_id=id"`
	Variant       string    `bun:"variant"`
	Description   string    `bun:"description"`
	OS            string    `bun:"os"`
	Arch          string    `bun:"arch"`
	ReleasedAt    time.Time `bun:"released_at"`
	Version       string    `bun:"version"`
	Unstable      bool      `bun:"unstable"`
	Alias         string    `bun:"alias"`
	Link          string    `bun:"link"`
	Format        string    `bun:"format"`
	Signature     string    `bun:"signature"` // a cryptographical representation (hash etc)
	Tags          []string  `bun:"tags,array"`
	UpgradeTarget string    `bun:"upgrade_target"` // If empty, the default upgrade target will be used
	ShouldUpgrade int       `bun:"should_upgrade"`
	UpdatedAt     time.Time `bun:"updated_at"`
	Groups        []Group   `bun:"m2m:releases_groups,join:Release=Group"`
}

func (r *Release) ToCatalogRelease() *catalog.Release {
	release := &catalog.Release{
		App: &catalog.App{
			Vendor:  r.App.Vendor,
			Product: r.App.Product,
			Name:    r.App.Name,
		},
		Variant:       r.Variant,
		Description:   r.Description,
		OS:            r.OS,
		Arch:          r.Arch,
		Date:          r.ReleasedAt,
		Version:       r.Version,
		Unstable:      r.Unstable,
		Alias:         r.Alias,
		Link:          r.Link,
		Format:        r.Format,
		Signature:     r.Signature,
		Tags:          r.Tags,
		UpgradeTarget: catalog.UpgradeTarget(r.UpgradeTarget),
		ShouldUpgrade: catalog.Criticality(r.ShouldUpgrade),
	}

	// Transform groups to simple strings
	release.Groups = make([]string, len(r.Groups))
	for i, group := range r.Groups {
		release.Groups[i] = group.Name
	}

	return release
}

func FromCatalogRelease(r *catalog.Release) Release {
	return Release{
		Variant:       r.Variant,
		Description:   r.Description,
		OS:            r.OS,
		Arch:          r.Arch,
		ReleasedAt:    r.Date,
		Version:       r.Version,
		Unstable:      r.Unstable,
		Alias:         r.Alias,
		Link:          r.Link,
		Format:        r.Format,
		Signature:     r.Signature,
		Tags:          r.Tags,
		UpgradeTarget: string(r.UpgradeTarget),
		ShouldUpgrade: int(r.ShouldUpgrade),
		UpdatedAt:     time.Now(),
	}
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

func (store *DbCatalogStore) FetchReleases(filter catalog.Filter, limit int) ([]*catalog.Release, error) {
	dbApp, err := store.getApp(filter.Vendor, filter.Product, false, false)
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

func (store *DbCatalogStore) StoreRelease(
	release *catalog.Release,
	allowUpdate bool,
) (*catalog.Release, error) {
	r := FromCatalogRelease(release)

	dbApp, err := store.getApp(release.App.Vendor, release.App.Product, true, true)
	if err != nil {
		return nil, err
	}

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

	// Check variant before inserting the release
	if r.Variant != "" {
		variantFound := false
		if dbApp.Variants != nil && len(dbApp.Variants) > 0 {
			for _, variant := range dbApp.Variants {
				if r.Variant == variant.Name {
					variantFound = true
					break
				}
			}
		}

		if !variantFound {
			return nil, fmt.Errorf("unknown variant \"%s\"", r.Variant)
		}
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
	if filter.Variant != "" || filter.EnforceVariant {
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
