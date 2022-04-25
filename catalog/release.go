package catalog

import (
	"time"
)

type Release struct {
	Vendor        string
	Product       string
	Name          string
	Variant       string
	Description   string
	OS            string
	Arch          string
	Date          time.Time
	Version       string // for use with semantic versioning
	Unstable      bool
	Alias         string
	Signature     string // a cryptographical representation (hash etc)
	Tags          []string
	UpgradeTarget UpgradeTarget // If empty, the default upgrade target will be used
	ShouldUpgrade Criticality
}

type ReleaseHistory struct {
	Vendor   string
	Product  string
	Releases []Release
}
