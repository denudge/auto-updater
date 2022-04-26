package catalog

type ClientState struct {
	Vendor         string
	Product        string
	Variant        string // optional
	OS             string // optional, e.g. for jars
	Arch           string // optional
	CurrentVersion string
	WithUnstable   bool
}

func (c ClientState) IsValid() bool {
	if c.Vendor == "" || c.Product == "" {
		return false
	}

	return true
}

func (c ClientState) IsInstalled() bool {
	if c.CurrentVersion == "" {
		return false
	}

	return true
}