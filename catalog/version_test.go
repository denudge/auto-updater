package catalog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindTargetVersion(t *testing.T) {
	versions := []string{
		"v1.0.0",
		"v1.0.1",
		"v1.1.0",
		"v1.1.1",
		"v1.2.0",
		"v1.2.1",
		"v1.2.2",
		"v1.3.0",
		"v1.4.0",
		"v1.5.0",
		"v1.7.0",
		"v1.7.1",
		"v1.8.0",
		"v1.8.1",
		"v2.0.0",
		"v2.0.1",
		"v2.0.2",
		"v2.1.0",
		"v2.2.0",
		"v2.3.0",
		"v2.3.1",
		"v2.3.2",
	}

	tcs := []struct {
		Name     string
		Current  string // current version
		Spec     string // UpgradeTarget
		Expected string
	}{
		{"Simple patch upgrade 1", "v2.3.1", "#.#.*", "v2.3.2"},
		{"Simple patch upgrade 2", "v1.1.0", "#.#.*", "v1.1.1"},
		{"Already at latest patch", "v1.1.1", "#.#.*", ""},
		{"Default upgrade: Minor", "v1.7.1", string(DefaultUpgradeTarget), "v1.8.1"},
		{"Default upgrade: Major", "v1.8.0", string(DefaultUpgradeTarget), "v2.0.2"},
		{"Default upgrade: None available", "v2.3.2", string(DefaultUpgradeTarget), ""},
		{"Next even minor", "v1.4.0", "#.+e._", "v1.8.0"},
		{"Next minor matching current even", "v1.2.0", "#.+c.*", "v1.4.0"},
		{"Next minor matching current odd", "v2.1.0", "#.+c.*", "v2.3.2"},
		{"Exact minor", "v1.3.0", "1.8.*", "v1.8.1"},
		{"Exact minor with Prefix", "v1.7.0", "v2.1.0", "v2.1.0"},
		{"Unavailable exact minor", "v1.5.0", "1.9.*", ""},
		{"2 minor steps", "v1.0.0", "#.+2.*", "v1.2.2"},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			result, err := FindTargetVersion(versions, tc.Current, UpgradeTarget(tc.Spec))
			assert.Nil(t, err)
			assert.Equal(t, tc.Expected, result)
		})
	}
}
