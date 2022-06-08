package catalog

import (
	"fmt"
	"time"
)

type Variant struct {
	App           *App
	Name          string
	Active        bool          // if this app variant is "delivered" or "handled" at all
	Locked        bool          // if updates of this app variant are "delivered" at all
	AllowRegister bool          // if clients are allowed to register
	UpgradeTarget UpgradeTarget // If empty, the upgrade target of the app will be used
	Created       time.Time
	Updated       time.Time
	DefaultGroups []string // empty means "default groups of the app" here
}

func (v *Variant) String() string {
	return fmt.Sprintf("%s %s", v.App.String(), v.Name)
}
