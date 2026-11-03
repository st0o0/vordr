package trivy

import (
	"encoding/json"
	"testing"
)

func TestParseTrivyOutput(t *testing.T) {
	raw := `{
		"SchemaVersion": 2,
		"ArtifactName": "nginx:latest",
		"Results": [{
			"Target": "nginx:latest (debian 12.7)",
			"Vulnerabilities": [
				{"Severity": "HIGH", "FixedVersion": "1.2.3"},
				{"Severity": "HIGH", "FixedVersion": ""},
				{"Severity": "MEDIUM", "FixedVersion": "2.0.0"},
				{"Severity": "LOW", "FixedVersion": ""},
				{"Severity": "", "FixedVersion": ""}
			]
		}]
	}`

	var report trivyReport
	if err := json.Unmarshal([]byte(raw), &report); err != nil {
		t.Fatalf("parse error: %v", err)
	}

	result := &ScanResult{
		Counts:  make(map[string]int),
		Fixable: make(map[string]int),
	}

	for _, r := range report.Results {
		for _, v := range r.Vulnerabilities {
			sev := v.Severity
			if sev == "" {
				sev = "UNKNOWN"
			}
			result.Counts[sev]++
			if v.FixedVersion != "" {
				result.Fixable[sev]++
			}
		}
	}

	if result.Counts["HIGH"] != 2 {
		t.Fatalf("expected HIGH=2, got %d", result.Counts["HIGH"])
	}
	if result.Counts["MEDIUM"] != 1 {
		t.Fatalf("expected MEDIUM=1, got %d", result.Counts["MEDIUM"])
	}
	if result.Counts["LOW"] != 1 {
		t.Fatalf("expected LOW=1, got %d", result.Counts["LOW"])
	}
	if result.Counts["UNKNOWN"] != 1 {
		t.Fatalf("expected UNKNOWN=1, got %d", result.Counts["UNKNOWN"])
	}
	if result.Fixable["HIGH"] != 1 {
		t.Fatalf("expected fixable HIGH=1, got %d", result.Fixable["HIGH"])
	}
	if result.Fixable["MEDIUM"] != 1 {
		t.Fatalf("expected fixable MEDIUM=1, got %d", result.Fixable["MEDIUM"])
	}
}

func TestParseTrivyEmptyResults(t *testing.T) {
	raw := `{"SchemaVersion": 2, "Results": [{"Vulnerabilities": null}]}`

	var report trivyReport
	if err := json.Unmarshal([]byte(raw), &report); err != nil {
		t.Fatalf("parse error: %v", err)
	}

	result := &ScanResult{
		Counts:  make(map[string]int),
		Fixable: make(map[string]int),
	}

	for _, r := range report.Results {
		for _, v := range r.Vulnerabilities {
			result.Counts[v.Severity]++
		}
	}

	if len(result.Counts) != 0 {
		t.Fatalf("expected empty counts, got %v", result.Counts)
	}
}

func TestParseTrivyMultipleResults(t *testing.T) {
	raw := `{
		"SchemaVersion": 2,
		"Results": [
			{"Vulnerabilities": [{"Severity": "CRITICAL", "FixedVersion": "1.0"}]},
			{"Vulnerabilities": [{"Severity": "CRITICAL", "FixedVersion": ""}, {"Severity": "LOW", "FixedVersion": ""}]}
		]
	}`

	var report trivyReport
	if err := json.Unmarshal([]byte(raw), &report); err != nil {
		t.Fatalf("parse error: %v", err)
	}

	result := &ScanResult{
		Counts:  make(map[string]int),
		Fixable: make(map[string]int),
	}

	for _, r := range report.Results {
		for _, v := range r.Vulnerabilities {
			sev := v.Severity
			if sev == "" {
				sev = "UNKNOWN"
			}
			result.Counts[sev]++
			if v.FixedVersion != "" {
				result.Fixable[sev]++
			}
		}
	}

	if result.Counts["CRITICAL"] != 2 {
		t.Fatalf("expected CRITICAL=2, got %d", result.Counts["CRITICAL"])
	}
	if result.Fixable["CRITICAL"] != 1 {
		t.Fatalf("expected fixable CRITICAL=1, got %d", result.Fixable["CRITICAL"])
	}
}
