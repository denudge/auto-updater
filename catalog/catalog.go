package catalog

type Catalog interface {
	ShouldUpgrade(OS string, Arch string, CurrentVersion string, WithUnstable bool) (Criticality, error)
	UpgradePath(OS string, Arch string, CurrentVersion string, WithUnstable bool) (UpgradePath, error)
}
