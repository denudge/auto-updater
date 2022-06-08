package api

import (
	"fmt"
	"github.com/denudge/auto-updater/catalog"
	"net/http"
)

type RegisterRequest struct {
	Vendor  string `json:"vendor"`
	Product string `json:"product"`
	Variant string `json:"variant"`
}

func (r *RegisterRequest) Validate() error {
	if r.Vendor == "" || r.Product == "" {
		return fmt.Errorf("vendor and product must be given")
	}

	// Add other rules here if necessary
	return nil
}

type RegisterResponse struct {
	ClientId string `json:"client_id"`
	Vendor   string `json:"vendor"`
	Product  string `json:"product"`
	Variant  string `json:"variant"`
}

func NewRegisterResponse(state *catalog.ClientState) *RegisterResponse {
	return &RegisterResponse{
		ClientId: state.ClientId,
		Vendor:   state.Vendor,
		Product:  state.Product,
		Variant:  state.Variant,
	}
}

func (api *Api) register(w http.ResponseWriter, r *http.Request) {
	request := RegisterRequest{}

	if err := api.parseAndValidateJsonPostRequest(w, r, &request); err != nil {
		// Errors are already written to the response
		return
	}

	clientState, err := api.app.RegisterClient(request.Vendor, request.Product, request.Variant)
	if err != nil {
		http.Error(w, fmt.Sprintf("error registering client: %s", err.Error()), http.StatusBadRequest)
		return
	}

	_ = api.writeJsonResponse(w, NewRegisterResponse(clientState))
}
