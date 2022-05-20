package database

import (
	"github.com/denudge/auto-updater/catalog"
	"github.com/uptrace/bun"
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
