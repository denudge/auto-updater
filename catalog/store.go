package catalog

type StoreInterface interface {
	Store(release *Release, allowUpdate bool) (*Release, error)
	Fetch(filter Filter) ([]*Release, error)
	SetCriticality(filter Filter, criticality Criticality) ([]*Release, error)
	SetStability(filter Filter, stability bool) ([]*Release, error)
	SetUpgradeTarget(filter Filter, target UpgradeTarget) ([]*Release, error)
}
