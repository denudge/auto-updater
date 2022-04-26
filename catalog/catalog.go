package catalog

type Catalog interface {
	ShouldUpgrade(state ClientState) (Criticality, error)
	// FindNextUpgrade should return nil if no update is available
	FindNextUpgrade(state ClientState) (*UpgradeStep, error)
	// FindUpgradePath should return nil if no update is available
	FindUpgradePath(state ClientState) (*UpgradePath, error)
}
