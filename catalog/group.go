package catalog

import "time"

type Group struct {
	App     *App
	Name    string
	Created time.Time
	Updated time.Time
}
