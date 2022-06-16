package app

import (
	"fmt"
	"github.com/denudge/auto-updater/cmd/updater/api"
	"github.com/denudge/auto-updater/updater"
	"strings"
)

type Updater struct {
	ConfigFileName string
	BaseUrl        string
	State          *updater.State
	Client         *api.CatalogClient
}

func NewUpdater(configFileName string, baseUrl string) *Updater {
	state, _ := updater.ReadStateFromFile(configFileName)

	return &Updater{
		ConfigFileName: configFileName,
		BaseUrl:        baseUrl,
		State:          state, // might be nil here
		Client:         api.NewApiClient(baseUrl),
	}
}

func (u *Updater) CheckStateConfiguration() error {
	if u.State == nil {
		return fmt.Errorf("configuration is missing. Please use \"updater init\" first")
	}

	return nil
}

func (u *Updater) Init(configFileName string, baseUrl string) {
	if configFileName == "" && baseUrl == "" {
		return
	}

	if configFileName != "" {
		u.ConfigFileName = configFileName

		state, _ := updater.ReadStateFromFile(u.ConfigFileName)
		u.State = state
	}

	if baseUrl != "" && baseUrl != u.BaseUrl {
		u.BaseUrl = baseUrl
		u.Client = api.NewApiClient(u.BaseUrl)
	}
}

func (u *Updater) CheckServerConfiguration() error {
	if u.BaseUrl != "" && !strings.HasPrefix(u.BaseUrl, "http") && u.Client != nil {
		return fmt.Errorf("catalog server address is missing or malformed (must start with \"http\")")
	}

	return nil
}

func (u *Updater) SaveState(state *updater.State) error {
	u.State = state

	return state.SaveToFile(u.ConfigFileName)
}
