package catalog

import "fmt"

type Criticality int

const (
	None                Criticality = 0
	Possible            Criticality = 1
	Recommended         Criticality = 2
	StronglyRecommended Criticality = 3
	Critical            Criticality = 4
	Enforced            Criticality = 5
	Exceptional         Criticality = -1 // manual interception necessary
)

func (c Criticality) String() string {
	switch c {
	case None:
		return "None"
	case Possible:
		return "Possible"
	case Recommended:
		return "Recommended"
	case StronglyRecommended:
		return "Strongly Recommended"
	case Critical:
		return "Critical"
	case Enforced:
		return "Enforced"
	case Exceptional:
		return "Exceptional"
	}

	return "Unknown"
}

func CriticalityFromString(value string) (Criticality, error) {
	switch value {
	case "None":
		return None, nil
	case "Possible":
		return Possible, nil
	case "Recommended":
		return Recommended, nil
	case "Strongly Recommended":
		return StronglyRecommended, nil
	case "Critical":
		return Critical, nil
	case "Enforced":
		return Enforced, nil
	case "Exceptional":
		return Exceptional, nil
	}

	return 0, fmt.Errorf("unknown criticality \"%s\"", value)
}
