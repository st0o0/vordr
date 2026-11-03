package scanner

import "time"

type Status int

const (
	Pending  Status = iota // 0 — discovered, waiting for scan
	Scanning               // 1 — scan in progress
	Scanned                // 2 — scan complete, results available
	Error                  // 3 — scan failed, will retry on next cycle
)

func (s Status) String() string {
	switch s {
	case Pending:
		return "pending"
	case Scanning:
		return "scanning"
	case Scanned:
		return "scanned"
	case Error:
		return "error"
	default:
		return "unknown"
	}
}

type ImageRef struct {
	Image   string // e.g. "grafana/grafana" — for Prometheus labels
	Tag     string // e.g. "13.3.0" — for Prometheus labels
	ImageID string // e.g. "sha256:52d09..." — for Trivy scan target
}

func (r ImageRef) String() string {
	return r.Image + ":" + r.Tag
}

func (r ImageRef) ScanTarget() string {
	if r.ImageID != "" {
		return r.ImageID
	}
	return r.String()
}

type Vulnerability struct {
	ID               string
	PkgName          string
	InstalledVersion string
	FixedVersion     string
	Severity         string
}

type ImageState struct {
	Ref             ImageRef
	Status          Status
	Containers      []string
	Counts          map[string]int // severity → count
	Fixable         map[string]int // severity → fixable count
	ScannedAt       time.Time
	Duration        time.Duration
	Error           string
	Vulnerabilities []Vulnerability
}
