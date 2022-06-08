package catalog

import (
	"strings"
	"time"
)

type Group struct {
	App     *App
	Name    string
	Created time.Time
	Updated time.Time
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
