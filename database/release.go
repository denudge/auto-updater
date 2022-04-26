package database

import (
	"github.com/denudge/auto-updater/catalog"
	"github.com/uptrace/bun"
	"time"
)

type Release struct {
	bun.BaseModel `bun:"table:releases"`
	ID            int64     `bun:"id,pk,autoincrement"`
	Vendor        string    `bun:"vendor"`
	Product       string    `bun:"product"`
	Name          string    `bun:"name"`
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
}

func (r *Release) ToCatalogRelease() *catalog.Release {
	return &catalog.Release{
		Vendor:        r.Vendor,
		Product:       r.Product,
		Name:          r.Name,
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
		UpgradeTarget: catalog.UpgradeTarget(r.UpgradeTarget), // If empty, the default upgrade target will be used
		ShouldUpgrade: catalog.Criticality(r.ShouldUpgrade),
	}
}
