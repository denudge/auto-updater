package catalog

import (
	"fmt"
	"time"
)

type Release struct {
	App           *App
	Variant       string // sth like "Pro" or "Light" (if given it must match a client's installation)
	Description   string
	OS            string // sth like "MacOS", "darwin" or "linux"; if given it must match a client's installation
	Arch          string // sth like "i386 or "ppc64"; if given it must match a client's installation
	Date          time.Time
	Version       string // for use with semantic versioning
	Unstable      bool   // some clients may not want to use unstable versions
	Alias         string // sth like "Focal Fossa"; optional, for printing (and release management) only
	Signature     string // a cryptographical representation (hash etc)
	Tags          []string
	UpgradeTarget UpgradeTarget // If empty, the default upgrade target will be used
	ShouldUpgrade Criticality
}

type ReleaseHistory struct {
	App      *App
	Releases []Release
}

func (r *Release) String() string {
	productName := r.App.String()

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

	version := r.Version
	if r.Alias != "" {
		version = fmt.Sprintf("%s (\"%s\")", r.Version, r.Alias)
	}

	if r.Unstable {
		version = fmt.Sprintf("%s [unstable]", version)
	}

	return fmt.Sprintf("%s%s, version %s, released on %s", productName, arch, version, r.Date.Format(time.RFC1123))
}
