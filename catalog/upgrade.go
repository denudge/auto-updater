package catalog

import (
	"fmt"
	"golang.org/x/mod/semver"
	"strings"
)

type UpgradeStep struct {
	Info        UpgradeInfo
	Release     Release
	Criticality Criticality
}

type UpgradePath struct {
	Info        UpgradeInfo
	Criticality Criticality
	Steps       []UpgradeStep
}

type UpgradeInfo struct {
	ShortInfo    string
	Description  string
	Explanation  string
	ReferenceUrl string
}

// FindInstallVersion retrieves the latest (stable) version
func FindInstallVersion(availableReleases []*Release, withUnstable bool) (*UpgradeStep, error) {
	if len(availableReleases) < 1 {
		return nil, nil
	}

	releaseMap, availableVersions := buildReleaseMap(availableReleases)

	// assume availableReleases to be presorted
	for i := len(availableVersions) - 1; i > 0; i-- {
		if !withUnstable && releaseMap[availableVersions[i]].Unstable {
			continue
		}

		return &UpgradeStep{
			Criticality: None,
			Release:     *releaseMap[availableVersions[i]],
		}, nil
	}

	// everything unstable or nothing released yet
	return nil, nil
}

func FindUpgradePath(availableReleases []*Release, currentVersion string) (*UpgradePath, error) {
	path := &UpgradePath{
		Steps: make([]UpgradeStep, 0, 4), // reserve spots for patch, minor and major
	}

	stepVersion := currentVersion
	for {
		step, err := FindNextUpgrade(availableReleases, stepVersion)
		if err != nil {
			return nil, err
		}

		if step == nil {
			// nothing new in town?
			if len(path.Steps) < 1 {
				return nil, nil
			}

			// summary = whatever the first steps dictates
			path.Criticality = path.Steps[0].Criticality
			path.Info = path.Steps[0].Info

			return path, nil
		}

		path.Steps = append(path.Steps, *step)
		stepVersion = step.Release.Version
	}
}

func FindNextUpgrade(availableReleases []*Release, currentVersion string) (*UpgradeStep, error) {
	releaseMap, availableVersions := buildReleaseMap(availableReleases)

	currentVersion = completeVersion(currentVersion)

	currentVersionObj, ok := releaseMap[currentVersion]
	if !ok {
		return nil, fmt.Errorf("current version not found in available releases")
	}

	target := currentVersionObj.UpgradeTarget

	// Check if we should search for patches if no further upgrade version has been found
	searchPatches := true
	if strings.HasPrefix(string(target), "nopatches:") {
		searchPatches = false
		target = target[len("nopatches:"):]
	}

	if target == "" {
		target = DefaultUpgradeTarget
	}

	targetVersion, err := FindTargetVersion(availableVersions, currentVersion, target)
	if err != nil {
		return nil, err
	}

	// nothing new in town
	if targetVersion == "" {
		if !searchPatches {
			return nil, nil
		}

		targetVersion, err = FindTargetVersion(availableVersions, currentVersion, "#.#.*")
		if err != nil {
			return nil, err
		}

		if targetVersion == "" {
			return nil, nil
		}
	}

	targetVersionObj, ok := releaseMap[targetVersion]
	if !ok {
		return nil, fmt.Errorf("target version not found in available releases")
	}

	criticality := calculateDefaultCriticality(currentVersion, targetVersion)
	if currentVersionObj.ShouldUpgrade > criticality {
		criticality = currentVersionObj.ShouldUpgrade
	}

	return &UpgradeStep{
		// TODO: Gather info and description from release store
		Release:     *targetVersionObj,
		Criticality: criticality,
	}, nil
}

func buildReleaseMap(releases []*Release) (releaseMap map[string]*Release, availableVersions []string) {
	// We can allocate the exact size needed here
	releaseMap = make(map[string]*Release, len(releases))
	availableVersions = make([]string, len(releases))

	for i, release := range releases {
		version := completeVersion(release.Version)

		availableVersions[i] = version
		releaseMap[version] = release
	}

	semver.Sort(availableVersions)

	return releaseMap, availableVersions
}

func calculateDefaultCriticality(currentVersion string, targetVersion string) Criticality {
	criticality := Possible
	if semver.MajorMinor(currentVersion) != semver.MajorMinor(targetVersion) {
		criticality = Recommended
	}
	if semver.Major(currentVersion) != semver.Major(targetVersion) {
		criticality = StronglyRecommended
	}

	return criticality
}

func (step *UpgradeStep) ToPath() *UpgradePath {
	return &UpgradePath{
		Steps:       []UpgradeStep{*step},
		Info:        step.Info,
		Criticality: step.Criticality,
	}
}
