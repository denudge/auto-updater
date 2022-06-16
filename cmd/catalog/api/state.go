package api

import (
	"fmt"
	"github.com/denudge/auto-updater/catalog"
)

type ClientStateRequest struct {
	ClientId       string `json:"client_id"`
	Vendor         string `json:"vendor"`
	Product        string `json:"product"`
	Variant        string `json:"variant"`
	OS             string `json:"os"`
	Arch           string `json:"arch"`
	WithUnstable   bool   `json:"with_unstable"`
	CurrentVersion string `json:"current_version"`
}

func NewClientStateRequest(state *catalog.ClientState) *ClientStateRequest {
	return &ClientStateRequest{
		ClientId:       state.ClientId,
		Vendor:         state.Vendor,
		Product:        state.Product,
		Variant:        state.Variant,
		OS:             state.OS,
		Arch:           state.Arch,
		WithUnstable:   state.WithUnstable,
		CurrentVersion: state.CurrentVersion,
	}
}

func (r *ClientStateRequest) ToClientState() *catalog.ClientState {
	return &catalog.ClientState{
		ClientId:       r.ClientId,
		Vendor:         r.Vendor,
		Product:        r.Product,
		Variant:        r.Variant,
		OS:             r.OS,
		Arch:           r.Arch,
		WithUnstable:   r.WithUnstable,
		CurrentVersion: r.CurrentVersion,
	}
}

func (r *ClientStateRequest) Validate() error {
	if r.Vendor == "" || r.Product == "" {
		return fmt.Errorf("vendor and product must be given")
	}

	// Add other rules here if necessary
	return nil
}
