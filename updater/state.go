package updater

import "time"

type State struct {
	Server      string
	ClientId    string
	Vendor      string
	Product     string
	Variant     string
	OS          string
	Arch        string
	Version     string
	LastChecked time.Time
}
