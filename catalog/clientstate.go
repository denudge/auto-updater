package catalog

type ClientState struct {
	ClientId       string
	Vendor         string
	Product        string
	Variant        string // optional
	OS             string // optional, e.g. for jars
	Arch           string // optional
	CurrentVersion string
	WithUnstable   bool
}

func (state ClientState) IsValid() bool {
	if state.ClientId == "" || state.Vendor == "" || state.Product == "" {
		return false
	}

	return true
}

func (state ClientState) IsInstalled() bool {
	if state.CurrentVersion == "" {
		return false
	}

	return true
}

func (state ClientState) ToFilter() Filter {
	filter := Filter{
		Vendor:         state.Vendor,
		Product:        state.Product,
		Variant:        state.Variant,
		EnforceVariant: true,
	}

	if state.IsInstalled() {
		filter.MinVersion = state.CurrentVersion
	}

	return filter
}