package catalog

import (
	"errors"
	"fmt"
	"golang.org/x/mod/semver"
	"regexp"
	"strconv"
	"strings"
)

type VersionInfo struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease string
	Build      string
}

func FindTargetVersion(availableVersions []string, currentVersion string, target UpgradeTarget) (string, error) {
	if !semver.IsValid(currentVersion) {
		return "", errors.New("current version is not a valid semantic version: " + currentVersion)
	}

	if !target.IsValid() {
		return "", errors.New("target string is not a valid upgrade target: " + string(target))
	}

	currentVersionInfo, err := parseVersionInfo(currentVersion)
	if err != nil {
		return "", err
	}

	semver.Sort(availableVersions)

	targets := strings.Split(string(target), ";")

	for _, targetStep := range targets {
		spec, err := UpgradeTarget(targetStep).FirstTargetSpec()
		if err != nil {
			return "", errors.New("could not extract first version spec from upgrage target: " + targetStep)
		}

		available, err := filterMajorVersion(availableVersions, spec.Major, currentVersionInfo.Major)
		if err != nil {
			return "", errors.New("could not filter major versions: %s" + err.Error())
		}

		if len(available) < 1 {
			continue
		}

		available, err = filterMinorVersion(available, spec.Minor, currentVersionInfo.Minor)
		if err != nil {
			return "", errors.New("could not filter minor versions: %s" + err.Error())
		}

		if len(available) < 1 {
			continue
		}

		available, err = filterPatchVersion(available, spec.Patch, currentVersionInfo.Patch)
		if err != nil {
			return "", errors.New("could not filter minor versions: %s" + err.Error())
		}

		if len(available) < 1 {
			continue
		}

		// Is there anything new here?
		availableVersion := available[len(available)-1]
		if semver.Compare(availableVersion, currentVersion) > 0 {
			return availableVersion, nil
		}
	}

	// nothing new in town
	return "", nil
}

func parseVersionInfo(version string) (VersionInfo, error) {
	if !semver.IsValid(version) {
		return VersionInfo{}, fmt.Errorf("is not a valid semantic version: %s", version)
	}

	major, err := extractMajor(version)
	if err != nil {
		return VersionInfo{}, err
	}

	minor, err := extractMinor(version)
	if err != nil {
		return VersionInfo{}, err
	}

	patch, err := extractPatch(version)
	if err != nil {
		return VersionInfo{}, err
	}

	return VersionInfo{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		Prerelease: semver.Prerelease(version),
		Build:      semver.Build(version),
	}, nil
}

func filterMajorVersion(availableVersions []string, spec string, currentValue int) ([]string, error) {
	if len(availableVersions) < 1 {
		return []string{}, nil
	}

	out := make([]string, 0, len(availableVersions))

	allMajors := extractAllMajors(availableVersions)

	targetVersion, err := findBySpecInSlice(allMajors, spec, currentValue)
	if err != nil {
		return []string{}, err
	}

	wanted := fmt.Sprintf("v%d", targetVersion)

	// assume availableVersions to be pre-sorted
	for _, version := range availableVersions {
		if semver.Major(version) == wanted {
			out = append(out, version)
		}
	}

	return out, nil
}

func filterMinorVersion(availableVersions []string, spec string, currentValue int) ([]string, error) {
	if len(availableVersions) < 1 {
		return []string{}, nil
	}

	out := make([]string, 0, len(availableVersions))

	allMinors := extractAllMinors(availableVersions)

	targetVersion, err := findBySpecInSlice(allMinors, spec, currentValue)
	if err != nil {
		return []string{}, err
	}

	// assume availableVersions to be pre-sorted
	for _, version := range availableVersions {
		minor, err := extractMinor(version)
		if err == nil && minor == targetVersion {
			out = append(out, version)
		}
	}

	return out, nil
}

func filterPatchVersion(availableVersions []string, spec string, currentValue int) ([]string, error) {
	if len(availableVersions) < 1 {
		return []string{}, nil
	}

	out := make([]string, 0, len(availableVersions))

	allPatches := extractAllPatches(availableVersions)

	targetVersion, err := findBySpecInSlice(allPatches, spec, currentValue)
	if err != nil {
		return []string{}, err
	}

	// assume availableVersions to be pre-sorted
	for _, version := range availableVersions {
		patch, err := extractPatch(version)
		if err == nil && patch == targetVersion {
			out = append(out, version)
		}
	}

	return out, nil
}

func findBySpecInSlice(slice []int, spec string, currentValue int) (int, error) {
	// TBD: Check if we should check existence in array or not...
	if spec == "*" {
		return latestOfSlice(slice)
	}

	if spec == "_" {
		return lowestOfSlice(slice)
	}

	if spec == "#" {
		return currentValue, nil
	}

	if spec[0] == '+' {
		nextSpec, err := parseNextSpec(spec, currentValue)
		if err != nil {
			return -1, err
		}

		return nextInSlice(slice, currentValue, nextSpec.Steps, nextSpec.Numbering)
	}

	if _, err := regexp.Match("v?\\d+", []byte(spec)); err == nil {
		numberStr := spec
		if numberStr[0] == 'v' {
			numberStr = numberStr[1:]
		}

		if exactNumber, err := strconv.Atoi(numberStr); err == nil {
			return exactNumber, nil
		}
	}

	return -1, fmt.Errorf("unknown spec format: %s", spec)
}

func latestOfSlice(slice []int) (int, error) {
	// If there is nothing more to find
	if len(slice) < 1 {
		return -1, errors.New("cannot find latest version in empty slice")
	}

	return slice[len(slice)-1], nil
}

func lowestOfSlice(slice []int) (int, error) {
	// If there is nothing more to find
	if len(slice) < 1 {
		return -1, errors.New("cannot find lowest version in empty slice")
	}

	return slice[0], nil
}

// Use only == "e" for even and "o" for odds
func nextInSlice(slice []int, currentValue int, steps int, only Numbering) (int, error) {
	// is -1 if not found (lets us start at zero position)
	currentOffset := sliceIndex(slice, currentValue)

	// already at the end?
	if currentOffset == len(slice)-1 {
		return -1, nil
	}

	for i := currentOffset + 1; i < len(slice) && steps > 0; i++ {
		pVersion := slice[i]
		if (only == OddNumbers && pVersion%2 == 1) ||
			(only == EvenNumbers && pVersion%2 == 0) ||
			(only == AnyNumber) {
			steps--
		}

		if steps == 0 {
			return pVersion, nil
		}
	}

	// We reached the end of all possible versions
	return -1, nil
}

func sliceIndex(slice []int, value int) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}

	return -1
}

func extractAllMajors(availableVersions []string) []int {
	out := make([]int, 0, len(availableVersions))

	// assume availableVersions to be pre-sorted
	for _, version := range availableVersions {
		major, err := extractMajor(version)
		if err == nil {
			if len(out) == 0 || out[len(out)-1] != major {
				out = append(out, major)
			}
		}
	}

	return out
}

func extractMajor(version string) (int, error) {
	canonical := semver.Canonical(version)
	major := semver.Major(canonical)

	return strconv.Atoi(strings.TrimPrefix(major, "v"))
}

func extractAllMinors(availableVersions []string) []int {
	out := make([]int, 0, len(availableVersions))

	// assume availableVersions to be pre-sorted
	for _, version := range availableVersions {
		minor, err := extractMinor(version)
		if err == nil {
			if len(out) == 0 || out[len(out)-1] != minor {
				out = append(out, minor)
			}
		}
	}

	return out
}

func extractMinor(version string) (int, error) {
	canonical := semver.Canonical(version)
	major := semver.Major(canonical)
	majorMinor := semver.MajorMinor(canonical)

	minor := majorMinor[len(major)+1:]

	return strconv.Atoi(minor)
}

func extractAllPatches(availableVersions []string) []int {
	out := make([]int, 0, len(availableVersions))

	// assume availableVersions to be pre-sorted
	for _, version := range availableVersions {
		patch, err := extractPatch(version)
		if err == nil {
			if len(out) == 0 || out[len(out)-1] != patch {
				out = append(out, patch)
			}
		}
	}

	return out
}

func extractPatch(version string) (int, error) {
	canonical := semver.Canonical(version)
	majorMinor := semver.MajorMinor(canonical)

	patch := canonical[len(majorMinor)+1:]

	if len(patch) < 1 {
		return -1, fmt.Errorf("could not extract patch from version: %s", version)
	}

	// extract leading numbers
	length := 0
	for ; length < len(patch); length++ {
		if patch[length] < '0' || patch[length] > '9' {
			break
		}
	}

	if length < 1 {
		return -1, fmt.Errorf("could not extract patch from version: %s", version)
	}

	return strconv.Atoi(patch[:length])
}
