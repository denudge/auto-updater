package catalog

import "golang.org/x/mod/semver"

type Filter struct {
	Vendor         string
	Product        string
	Name           string
	Variant        string
	EnforceVariant bool // if an empty variant will be matched
	OS             string
	Arch           string
	Alias          string
	MinVersion     string // use MinVersion == MaxVersion to hit an exact version
	MaxVersion     string
	AfterVersion   string // like a "MinVersion" but excluding this version
	BeforeVersion  string // like a "MaxVersion" but excluding this version
	WithUnstable   bool
	Groups         []string
}

func (f *Filter) CompleteVersions() {
	if f.MinVersion != "" {
		f.MinVersion = completeVersion(f.MinVersion)
	}

	if f.AfterVersion != "" {
		f.AfterVersion = completeVersion(f.AfterVersion)
	}

	if f.BeforeVersion != "" {
		f.BeforeVersion = completeVersion(f.BeforeVersion)
	}

	if f.MaxVersion != "" {
		f.MaxVersion = completeVersion(f.MaxVersion)
	}
}

func (f *Filter) MatchVersion(version string) bool {
	version = completeVersion(version)

	if f.MinVersion != "" && semver.Compare(f.MinVersion, version) > 0 {
		return false
	}

	if f.AfterVersion != "" && semver.Compare(f.AfterVersion, version) >= 0 {
		return false
	}

	if f.BeforeVersion != "" && semver.Compare(f.BeforeVersion, version) <= 0 {
		return false
	}

	if f.MaxVersion != "" && semver.Compare(f.MaxVersion, version) < 0 {
		return false
	}

	return true
}

func (f *Filter) FiltersVersions() bool {
	return f.MinVersion != "" || f.AfterVersion != "" || f.BeforeVersion != "" || f.MaxVersion != ""
}

type GroupFilter struct {
	Vendor  string
	Product string
	Name    string
}
