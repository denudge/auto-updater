package updater

import (
	"fmt"
	"github.com/denudge/auto-updater/catalog"
	"os"
	"strings"
	"time"
)

type State struct {
	Server       string
	ClientId     string
	Vendor       string
	Product      string
	Variant      string
	OS           string
	Arch         string
	WithUnstable bool
	Version      string
	LastChecked  time.Time
}

func (state *State) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil || file == nil {
		return fmt.Errorf("could not write state to file \"%s\": %s", filename, err.Error())
	}

	str := "server-address=" + state.Server + "\n" +
		"client-id=" + state.ClientId + "\n" +
		"vendor=" + state.Vendor + "\n" +
		"product=" + state.Product + "\n" +
		"variant=" + state.Variant + "\n" +
		"os=" + state.OS + "\n" +
		"arch=" + state.Arch + "\n" +
		"with-unstable=" + fmt.Sprintf("%v", state.WithUnstable) + "\n" +
		"version=" + state.Version + "\n" +
		"last-checked=" + state.LastChecked.Format(time.RFC3339) + "\n"

	if _, err = file.Write([]byte(str)); err != nil {
		return err
	}

	return nil
}

func (state *State) IsValid() bool {
	if state.Server == "" || state.Vendor == "" || state.Product == "" {
		return false
	}

	return true
}

func (state *State) IsInstalled() bool {
	return state.Version != ""
}

func (state *State) ToClientState() *catalog.ClientState {
	return &catalog.ClientState{
		ClientId:       state.ClientId,
		Vendor:         state.Vendor,
		Product:        state.Product,
		Variant:        state.Variant,
		CurrentVersion: state.Version,
		OS:             state.OS,
		Arch:           state.Arch,
		WithUnstable:   state.WithUnstable,
	}
}

func StateFromClientState(state *catalog.ClientState, server string, lastChecked time.Time) *State {
	return &State{
		Server:       server,
		ClientId:     state.ClientId,
		Vendor:       state.Vendor,
		Product:      state.Product,
		Variant:      state.Variant,
		Version:      state.CurrentVersion,
		OS:           state.OS,
		Arch:         state.Arch,
		WithUnstable: state.WithUnstable,
		LastChecked:  lastChecked,
	}
}

func ReadStateFromFile(filename string) (*State, error) {
	content, err := os.ReadFile(filename)
	if err != nil || len(content) < 1 {
		return nil, fmt.Errorf("could not read state from file \"%s\"", filename)
	}

	state := State{}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		before, after, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		switch before {
		case "server-address":
			state.Server = after
			break
		case "client-id":
			state.ClientId = after
			break
		case "vendor":
			state.Vendor = after
			break
		case "product":
			state.Product = after
			break
		case "variant":
			state.Variant = after
			break
		case "os":
			state.OS = after
			break
		case "arch":
			state.Arch = after
			break
		case "version":
			state.Version = after
			break
		case "with-unstable":
			if after == "yes" || after == "true" {
				state.WithUnstable = true
			}
			break
		case "last-checked":
			if after != "" {
				state.LastChecked, err = time.Parse(time.RFC3339, after)
				if err != nil {
					return nil, fmt.Errorf("invalid last-checked time format")
				}
			}
		default:
			return nil, fmt.Errorf("invalid state file format")
		}
	}

	return &state, nil
}
