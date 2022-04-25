package catalog

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
