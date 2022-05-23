package catalog

type StoreInterface interface {
	// App management
	ListApps(limit int) ([]*App, error)
	StoreApp(app *App, allowUpdate bool) (*App, error)
	FindApp(vendor string, product string) (*App, error)
	SetAppDefaultGroups(app *App) (*App, error)

	// Group management
	ListGroups(filter GroupFilter, limit int) ([]*Group, error)
	StoreGroup(group *Group, allowUpdate bool) (*Group, error)

	// Release management
	StoreRelease(release *Release, allowUpdate bool) (*Release, error)
	FetchReleases(filter Filter) ([]*Release, error)
	SetCriticality(filter Filter, criticality Criticality) ([]*Release, error)
	SetStability(filter Filter, stability bool) ([]*Release, error)
	SetUpgradeTarget(filter Filter, target UpgradeTarget) ([]*Release, error)
}
