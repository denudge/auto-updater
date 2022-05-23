package catalog

import (
	"fmt"
	"strings"
	"time"
)

type App struct {
	Vendor        string        // must be present and match a client's installation
	Product       string        // must be present and match a client's installation
	Name          string        // for printing; if not given, "<vendor> <product>" will be used
	Active        bool          // if this app is "delivered" or "handled" at all
	Locked        bool          // if updates of this app are "delivered" at all
	UpgradeTarget UpgradeTarget // If empty, the default upgrade target will be used
	Created       time.Time
	Updated       time.Time
	Groups        []string
	DefaultGroups []string // empty means "public" by default
}

func (app *App) String() string {
	if app.Name != "" {
		return app.Name
	}

	return fmt.Sprintf("%s %s", app.Vendor, app.Product)
}

func FormatGroups(groups []string) string {
	str := "("
	if len(groups) > 0 {
		str += strings.Join(groups, ", ")
	} else {
		str += "public"
	}
	str += ")"

	return str
}
