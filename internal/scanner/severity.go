package scanner

var severityOrder = map[string]int{
	"UNKNOWN": 0, "LOW": 1, "MEDIUM": 2, "HIGH": 3, "CRITICAL": 4,
}

func MeetsThreshold(severity, threshold string) bool {
	return severityOrder[severity] >= severityOrder[threshold]
}
