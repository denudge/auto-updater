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
	return FormatStringList(groups, "public")
}

func FormatVariants(variants []string) string {
	return FormatStringList(variants, "")
}

func FormatStringList(strs []string, empty string) string {
	if strs != nil && len(strs) > 0 {
		return "(" + strings.Join(strs, ", ") + ")"
	} else {
		return "(" + empty + ")"
	}
}
