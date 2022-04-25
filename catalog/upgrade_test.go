package catalog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindNextUpgrade(t *testing.T) {
	versions := []*Release{
		{Version: "v1.0.0", UpgradeTarget: "#.+2.*"},
		{Version: "v1.0.1"},
		{Version: "v1.1.0", UpgradeTarget: "#.#.*"},
		{Version: "v1.1.1", UpgradeTarget: "#.#.*"},
		{Version: "v1.2.0", UpgradeTarget: "#.+c.*"},
		{Version: "v1.2.1"},
		{Version: "v1.2.2"},
		{Version: "v1.3.0", UpgradeTarget: "1.8.*"},
		{Version: "v1.4.0", UpgradeTarget: "#.+e._"},
		{Version: "v1.5.0", UpgradeTarget: "1.9.*"},
		{Version: "v1.7.0", UpgradeTarget: "v2.1.0"},
		{Version: "v1.7.1"},
		{Version: "v1.8.0"},
		{Version: "v1.8.1"},
		{Version: "v2.0.0"},
		{Version: "v2.0.1"},
		{Version: "v2.0.2"},
		{Version: "v2.1.0", UpgradeTarget: "#.+c.*"},
		{Version: "v2.2.0"},
		{Version: "v2.3.0"},
		{Version: "v2.3.1", UpgradeTarget: "#.#.*"},
		{Version: "v2.3.2"},
	}

	tcs := []struct {
		Name     string
		Current  string // current version
		Expected string
	}{
		{"Simple patch upgrade 1", "v2.3.1", "v2.3.2"},
		{"Simple patch upgrade 2", "v1.1.0", "v1.1.1"},
		{"Already at latest patch", "v1.1.1", ""},
		{"Default upgrade: Minor", "v1.7.1", "v1.8.1"},
		{"Default upgrade: Major", "v1.8.0", "v2.0.2"},
		{"Default upgrade: None available", "v2.3.2", ""},
		{"Next even minor", "v1.4.0", "v1.8.0"},
		{"Next minor matching current even", "v1.2.0", "v1.4.0"},
		{"Next minor matching current odd", "v2.1.0", "v2.3.2"},
		{"Exact minor", "v1.3.0", "v1.8.1"},
		{"Unavailable exact minor", "v1.5.0", ""},
		{"2 minor steps", "v1.0.0", "v1.2.2"},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			result, err := FindNextUpgrade(versions, tc.Current)
			assert.Nil(t, err)

			resultStr := ""
			if result != nil {
				resultStr = result.Release.Version
			}
			assert.Equal(t, tc.Expected, resultStr)
		})
	}
}

func TestFindNextUpgradeWorksWithoutLeadingV(t *testing.T) {
	versions := []*Release{
		{Version: "1.1.0"},
		{Version: "1.1.1"},
	}

	tcs := []struct {
		Name     string
		Current  string // current version
		Expected string
	}{
		{"Simple patch upgrade", "1.1.0", "1.1.1"},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			result, err := FindNextUpgrade(versions, tc.Current)
			assert.Nil(t, err)

			resultStr := ""
			if result != nil {
				resultStr = result.Release.Version
			}
			assert.Equal(t, tc.Expected, resultStr)
		})
	}
}

func TestFindNextUpgradeIgnoresPatchesIfConfigured(t *testing.T) {
	versions := []*Release{
		{Version: "1.1.0", UpgradeTarget: "nopatches:"},
		{Version: "1.1.1"},
	}

	tcs := []struct {
		Name     string
		Current  string // current version
		Expected string
	}{
		{"Simple patch upgrade", "1.1.0", ""},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			result, err := FindNextUpgrade(versions, tc.Current)
			assert.Nil(t, err)

			resultStr := ""
			if result != nil {
				resultStr = result.Release.Version
			}
			assert.Equal(t, tc.Expected, resultStr)
		})
	}
}

func TestFindUpgradePath(t *testing.T) {
	versions := []*Release{
		{Version: "v1.0.0"},
		{Version: "v1.0.1"},
		{Version: "v1.1.0", UpgradeTarget: "#.#.*"},
		{Version: "v1.1.1", UpgradeTarget: "#.#.*"},
		{Version: "v1.2.0", UpgradeTarget: "#.+c.*"},
		{Version: "v1.2.1"},
		{Version: "v1.2.2"},
		{Version: "v1.3.0", UpgradeTarget: "1.8.*"},
		{Version: "v1.4.0", UpgradeTarget: "#.+e._"},
		{Version: "v1.5.0", UpgradeTarget: "1.9.*"},
		{Version: "v1.7.0", UpgradeTarget: "v2.1.0"},
		{Version: "v1.7.1"},
		{Version: "v1.8.0"},
		{Version: "v1.8.1"},
		{Version: "v2.0.0"},
		{Version: "v2.0.1"},
		{Version: "v2.0.2"},
		{Version: "v2.1.0", UpgradeTarget: "#.+c.*"},
		{Version: "v2.2.0"},
		{Version: "v2.3.0"},
		{Version: "v2.3.1", UpgradeTarget: "#.#.*"},
		{Version: "v2.3.2"},
	}

	tcs := []struct {
		Name           string
		Current        string // current version
		Expected       []string
		ExpCriticality Criticality
	}{
		{"Simple patch upgrade 1", "v2.3.1", []string{"v2.3.2"}, Possible},
		{"Simple patch upgrade 2", "v1.1.0", []string{"v1.1.1"}, Possible},
		{"Already at latest patch", "v1.1.1", []string{}, None},
		{"Default upgrade: Minor", "v1.7.1", []string{"v1.8.1", "v2.0.2", "v2.1.0", "v2.3.2"}, Recommended},
		{"Default upgrade: None available", "v2.3.2", []string{}, None},
		{"Next even minor", "v1.4.0", []string{"v1.8.0", "v2.0.2", "v2.1.0", "v2.3.2"}, Recommended},
		{"Next minor matching current odd", "v2.1.0", []string{"v2.3.2"}, Recommended},
		{"Unavailable exact minor", "v1.5.0", []string{}, None},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			result, err := FindUpgradePath(versions, tc.Current)
			assert.Nil(t, err)

			versionList := make([]string, 0, 4)
			criticality := None
			if result != nil {
				criticality = result.Criticality
				for _, step := range result.Steps {
					versionList = append(versionList, step.Release.Version)
				}
			}

			assert.Equal(t, tc.Expected, versionList)
			assert.Equal(t, tc.ExpCriticality, criticality)
		})
	}
}

func TestFindInstallVersion(t *testing.T) {
	versions := []*Release{
		{Version: "v1.0.0", Unstable: true},
		{Version: "v1.0.1", Unstable: true},
		{Version: "v1.1.0"},
		{Version: "v1.1.1"},
		{Version: "v1.2.0"},
		{Version: "v1.2.1"},
		{Version: "v1.2.2", Unstable: true},
	}

	result, err := FindInstallVersion(versions, false)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "v1.2.1", result.Release.Version)

	result, err = FindInstallVersion(versions, true)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "v1.2.2", result.Release.Version)

	versions = []*Release{
		{Version: "v1.0.0", Unstable: true},
		{Version: "v1.0.1", Unstable: true},
		{Version: "v1.1.0", Unstable: true},
	}

	result, err = FindInstallVersion(versions, false)
	assert.Nil(t, err)
	assert.Nil(t, result)
}
