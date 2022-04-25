package catalog

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

/*
UpgradeTarget defines the next logical version from a given, mostly current, semantic version
- Use a colon ";" to divide upgrade steps that must be run through sequentially
- Use a hashtag "#" to refer to the current number
- Use a star sign "*" to refer to the highest and latest available number
- A plus sign "+([1-9]|e|o|c)?" refers to the next available number.
  Using an optional step size will set the behaviour to skip (n-1) intermediate numbers
  (meaning current "7" out of a release history of [6,7,9,10] and "+2" will result in "10" as upgrade target).
  Use +e for the next even number and +o for the next odd number.
  Use +c for the next number that is even or odd depending on the current number.
- Use an underscore sign "_" to refer to the lowest available number.
- Use a number to refer to a specific version.

Examples:
- "*.*.*" refers to the absolute latest version.
- "#.#.*" refers to the latest patch version.
- "+._.*" refers to the latest patch of the first minor version of the next major version.
- "4.0.0" refers to 4.0.0 exactly

If no version is found matching the above spec, then at least the latest patch version is searched for.
If you want to exclude that, add a "nopatches:" prefix to the upgrade target string, e.g. "nopatches:*.*.0"
*/
type UpgradeTarget string

const (
	nextSpecPattern          = "\\+((\\d+)?(e|o|c)?)"
	modifierPattern          = "(nopatches\\:)?"
	prefixPattern            = "([v])?"
	numberPattern            = "([v]?\\d+|\\_|\\#|\\*|\\+(\\d+)?(e|o|c)?)"
	tripleNumberPattern      = numberPattern + "\\." + numberPattern + "\\." + numberPattern
	preReleasePattern        = "(\\-[^\\+]+)?"
	buildPattern             = "(\\+.*)?"
	upgradeTargetStepPattern = modifierPattern + prefixPattern + tripleNumberPattern + preReleasePattern + buildPattern
)

var nextSpecRegExp *regexp.Regexp
var upgradeTargetStepRegExp *regexp.Regexp

func init() {
	nextSpecRegExp = regexp.MustCompile(nextSpecPattern)
	upgradeTargetStepRegExp = regexp.MustCompile(upgradeTargetStepPattern)
}

func (target UpgradeTarget) IsValid() bool {
	targets := strings.Split(string(target), ";")

	for _, targetStep := range targets {
		if !upgradeTargetStepRegExp.Match([]byte(targetStep)) {
			return false
		}
	}

	return true
}

type TargetSpec struct {
	Major      string
	Minor      string
	Patch      string
	PreRelease string
	Build      string
}

type Numbering string

const (
	AnyNumber   Numbering = ""
	EvenNumbers           = "e"
	OddNumbers            = "o"
)

type NextSpec struct {
	Steps     int
	Numbering Numbering
}

func parseNextSpec(spec string, currentValue int) (NextSpec, error) {
	if len(spec) < 1 {
		return NextSpec{}, errors.New("cannot extract next spec from empty string")
	}

	if spec[0] != '+' {
		return NextSpec{}, errors.New("next spec has to start with plus sign (\"+\")")
	}

	parts := nextSpecRegExp.FindAllStringSubmatch(spec, 1)
	if parts == nil {
		return NextSpec{}, fmt.Errorf("next spec pattern mismatch: %s", spec)
	}

	if len(parts) < 1 || len(parts[0]) < 3 {
		return NextSpec{}, fmt.Errorf("next spec pattern mismatch: %s", spec)
	}

	steps := 1
	if len(parts[0][2]) > 0 {
		pSteps, err := strconv.Atoi(parts[0][2])
		if err != nil {
			return NextSpec{}, fmt.Errorf("next spec pattern mismatch: %s", spec)
		}

		steps = pSteps
	}

	numbering := AnyNumber
	if len(parts[0][3]) > 0 {
		if parts[0][3] == OddNumbers || parts[0][3] == EvenNumbers {
			numbering = Numbering(parts[0][3])
		}

		if parts[0][3] == "c" {
			if currentValue%2 == 1 {
				numbering = OddNumbers
			} else {
				numbering = EvenNumbers
			}
		}
	}

	return NextSpec{
		Steps:     steps,
		Numbering: numbering,
	}, nil
}

func (target UpgradeTarget) FirstTargetSpec() (*TargetSpec, error) {
	parts := upgradeTargetStepRegExp.FindAllStringSubmatch(string(target), 1)
	if parts == nil {
		return nil, errors.New("target string is not a valid upgrade target: " + string(target))
	}

	return &TargetSpec{
		Major:      parts[0][3],
		Minor:      parts[0][6],
		Patch:      parts[0][9],
		PreRelease: parts[0][12],
		Build:      parts[0][13],
	}, nil
}

const (
	// DefaultUpgradeTarget means: Upgrade to the latest minor+patch of the current major, then use the first minor of the next major
	DefaultUpgradeTarget UpgradeTarget = "#.+.*;+._.*"
	AlwaysLatest                       = "*.*.*" // Useful for uncritical, possibly stateless software
	OnlyMinorAndPatches                = "#.*.*" // Useful for end-of-life versions
	OnlyPatches                        = "#.#.*" // Useful for ensuring highest compatibility
)
