package app

import (
	"github.com/denudge/auto-updater/cmd/updater/api"
)

type Updater struct {
	Client *api.CatalogClient
}

func NewUpdater(baseUrl string) *Updater {
	return &Updater{
		Client: api.NewApiClient(baseUrl),
	}
}
