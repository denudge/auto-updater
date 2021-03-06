package api

import (
	"github.com/denudge/auto-updater/catalog"
	"github.com/denudge/auto-updater/cmd/catalog/api"
	"net/http"
	"time"
)

type CatalogClient struct {
	baseUrl string
	client  *http.Client
}

func NewApiClient(baseUrl string) *CatalogClient {
	return &CatalogClient{
		baseUrl: baseUrl,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (client *CatalogClient) RegisterClient(vendor string, product string, variant string) (*catalog.ClientState, error) {
	request := api.RegisterRequest{
		Vendor:  vendor,
		Product: product,
		Variant: variant,
	}

	response := api.RegisterResponse{}

	err := client.doJson(request, http.MethodPost, "/register", &response)
	if err != nil {
		return nil, err
	}

	return response.ToClientState(), nil
}

func (client *CatalogClient) FindNextUpgrade(state *catalog.ClientState) (*catalog.UpgradeStep, error) {
	request := api.NewClientStateRequest(state)

	response := api.UpgradeStepResponse{}

	err := client.doJson(request, http.MethodPost, "/upgrade/step", &response)
	if err != nil {
		return nil, err
	}

	step, err := response.ToUpgradeStep()
	if err != nil {
		return nil, err
	}

	return &step, nil
}
