package catalog

type StoreInterface interface {
	CreateApp(app *App, allowUpdate bool) (*App, error)
	StoreRelease(release *Release, allowUpdate bool) (*Release, error)
	FetchReleases(filter Filter) ([]*Release, error)
	SetCriticality(filter Filter, criticality Criticality) ([]*Release, error)
	SetStability(filter Filter, stability bool) ([]*Release, error)
	SetUpgradeTarget(filter Filter, target UpgradeTarget) ([]*Release, error)
}
