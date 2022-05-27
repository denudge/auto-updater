package catalog

import "time"

type Client struct {
	App     *App
	Variant string
	Uuid    string
	Name    string
	Active  bool
	Locked  bool
	Groups  []string // will be hydrated by DB layer
	Created time.Time
	Updated time.Time
}
