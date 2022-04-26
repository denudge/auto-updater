package catalog

type StoreInterface interface {
	Store(release *Release) error
	Fetch(filter Filter) ([]*Release, error)
	SetCriticality(filter Filter, criticality Criticality) ([]*Release, error)
	SetStability(filter Filter, stability bool) ([]*Release, error)
	SetUpgradeTarget(filter Filter, target UpgradeTarget) ([]*Release, error)
}

type Filter struct {
	Vendor        string
	Product       string
	Name          string
	Variant       string
	OS            string
	Arch          string
	Alias         string
	MinVersion    string // use MinVersion == MaxVersion to hit an exact version
	MaxVersion    string
	AfterVersion  string // like a "MinVersion" but excluding this version
	BeforeVersion string // like a "MaxVersion" but excluding this version
	WithUnstable  bool
}
