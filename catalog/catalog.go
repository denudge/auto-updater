package catalog

type Catalog interface {
	RegisterClient(vendor string, product string, variant string) (*ClientState, error)
	ShouldUpgrade(state *ClientState) (Criticality, error)
	// FindNextUpgrade should return nil if no update is available
	FindNextUpgrade(state *ClientState) (*UpgradeStep, error)
	// FindUpgradePath should return nil if no update is available
	FindUpgradePath(state *ClientState) (*UpgradePath, error)
}
