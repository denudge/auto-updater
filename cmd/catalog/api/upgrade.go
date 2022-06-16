package api

import (
	"fmt"
	"github.com/denudge/auto-updater/catalog"
	"net/http"
	"time"
)

type UpgradeInfoResponse struct {
	ShortInfo    string `json:"short_info,omitempty"`
	Description  string `json:"description,omitempty"`
	Explanation  string `json:"explanation,omitempty"`
	ReferenceUrl string `json:"reference_url,omitempty"`
}

func NewUpgradeInfoResponse(info catalog.UpgradeInfo) UpgradeInfoResponse {
	return UpgradeInfoResponse{
		ShortInfo:    info.ShortInfo,
		Description:  info.Description,
		Explanation:  info.Explanation,
		ReferenceUrl: info.ReferenceUrl,
	}
}

func (r *UpgradeInfoResponse) ToUpgradeInfo() catalog.UpgradeInfo {
	return catalog.UpgradeInfo{
		ShortInfo:    r.ShortInfo,
		Description:  r.Description,
		Explanation:  r.Explanation,
		ReferenceUrl: r.ReferenceUrl,
	}
}

type ReleaseResponse struct {
	Vendor      string   `json:"vendor"`
	Product     string   `json:"product"`
	Variant     string   `json:"variant"`
	Description string   `json:"description,omitempty"`
	OS          string   `json:"os,omitempty"`
	Arch        string   `json:"arch,omitempty"`
	Date        string   `json:"date"`
	Version     string   `json:"version"`
	Unstable    bool     `json:"unstable,omitempty"`
	Alias       string   `json:"alias,omitempty"`
	Link        string   `json:"link,omitempty"`
	Format      string   `json:"format,omitempty"`
	Signature   string   `json:"signature,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Criticality string   `json:"criticality"`
}

func (r *ReleaseResponse) ToRelease() (catalog.Release, error) {
	criticality, err := catalog.CriticalityFromString(r.Criticality)
	if err != nil {
		return catalog.Release{}, err
	}

	date, err := time.Parse(time.RFC3339, r.Date)
	if err != nil {
		return catalog.Release{}, err
	}

	return catalog.Release{
		App: &catalog.App{
			Vendor:  r.Vendor,
			Product: r.Product,
		},
		Variant:       r.Variant,
		Description:   r.Description,
		OS:            r.OS,
		Arch:          r.Arch,
		Date:          date,
		Version:       r.Version,
		Unstable:      r.Unstable,
		Alias:         r.Alias,
		Link:          r.Link,
		Format:        r.Format,
		Signature:     r.Signature,
		Tags:          r.Tags,
		ShouldUpgrade: criticality,
	}, nil
}

func NewReleaseResponse(r catalog.Release) ReleaseResponse {
	return ReleaseResponse{
		Vendor:      r.App.Vendor,
		Product:     r.App.Product,
		Variant:     r.Variant,
		Description: r.Description,
		OS:          r.OS,
		Arch:        r.Arch,
		Date:        r.Date.Format(time.RFC3339),
		Version:     r.Version,
		Unstable:    r.Unstable,
		Alias:       r.Alias,
		Link:        r.Link,
		Format:      r.Format,
		Signature:   r.Signature,
		Tags:        r.Tags,
		Criticality: r.ShouldUpgrade.String(),
	}
}

type UpgradeStepResponse struct {
	Info        UpgradeInfoResponse `json:"info,omitempty"`
	Release     ReleaseResponse     `json:"release"`
	Criticality string              `json:"criticality"`
}

func (u *UpgradeStepResponse) ToUpgradeStep() (catalog.UpgradeStep, error) {
	criticality, err := catalog.CriticalityFromString(u.Criticality)
	if err != nil {
		return catalog.UpgradeStep{}, err
	}

	release, err := u.Release.ToRelease()
	if err != nil {
		return catalog.UpgradeStep{}, err
	}

	return catalog.UpgradeStep{
		Info:        u.Info.ToUpgradeInfo(),
		Release:     release,
		Criticality: criticality,
	}, nil
}

func NewUpgradeStepResponse(step *catalog.UpgradeStep) UpgradeStepResponse {
	return UpgradeStepResponse{
		Info:        NewUpgradeInfoResponse(step.Info),
		Release:     NewReleaseResponse(step.Release),
		Criticality: step.Criticality.String(),
	}
}

func (api *Api) findNextUpgrade(w http.ResponseWriter, r *http.Request) {
	request := ClientStateRequest{}

	if err := api.parseAndValidateJsonPostRequest(w, r, &request); err != nil {
		// Errors are already written to the response
		return
	}

	upgradeStep, err := api.catalog.FindNextUpgrade(request.ToClientState())
	if err != nil {
		http.Error(w, fmt.Sprintf("error finding next upgrade: %s", err.Error()), http.StatusBadRequest)
		return
	}

	_ = api.writeJsonResponse(w, NewUpgradeStepResponse(upgradeStep))
}
