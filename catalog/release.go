package catalog

import (
	"fmt"
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

func (r *Release) String() string {
	productName := fmt.Sprintf("%s %s", r.Vendor, r.Product)
	if r.Variant != "" {
		productName = fmt.Sprintf("%s %s", productName, r.Variant)
	}

	arch := ""
	if r.OS != "" || r.Arch != "" {
		sep := ""
		if r.OS != "" && r.Arch != "" {
			sep = ", "
		}
		arch = fmt.Sprintf(" (%s%s%s)", r.OS, sep, r.Arch)
	}

	return fmt.Sprintf("%s%s %s, released on %s", productName, arch, r.Version, r.Date)
}
