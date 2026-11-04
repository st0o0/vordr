package trivy

import (
	"encoding/json"
	"testing"
)

func TestNewCLIStandaloneArgs(t *testing.T) {
	cli := NewCLI("/cache", 5*60000000000, "")
	if cli.serverURL != "" {
		t.Errorf("expected empty serverURL, got %q", cli.serverURL)
	}
	if cli.cacheDir != "/cache" {
		t.Errorf("expected cacheDir /cache, got %q", cli.cacheDir)
	}
}

func TestNewCLIServerArgs(t *testing.T) {
	cli := NewCLI("/cache", 5*60000000000, "http://trivy:4954")
	if cli.serverURL != "http://trivy:4954" {
		t.Errorf("expected serverURL http://trivy:4954, got %q", cli.serverURL)
	}
}

func TestParseTrivyOutput(t *testing.T) {
	raw := `{
		"SchemaVersion": 2,
		"ArtifactName": "nginx:latest",
		"Results": [{
			"Target": "nginx:latest (debian 12.7)",
			"Vulnerabilities": [
				{"VulnerabilityID": "CVE-2024-0001", "PkgName": "libfoo", "InstalledVersion": "1.0.0", "Severity": "HIGH", "FixedVersion": "1.2.3"},
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
			result.Vulnerabilities = append(result.Vulnerabilities, Vulnerability{
				ID:               v.VulnerabilityID,
				PkgName:          v.PkgName,
				InstalledVersion: v.InstalledVersion,
				FixedVersion:     v.FixedVersion,
				Severity:         sev,
			})
		}
	}

	if len(result.Vulnerabilities) != 5 {
		t.Fatalf("expected 5 vulnerabilities, got %d", len(result.Vulnerabilities))
	}
	if got, want := result.Vulnerabilities[0], (Vulnerability{
		ID: "CVE-2024-0001", PkgName: "libfoo", InstalledVersion: "1.0.0", FixedVersion: "1.2.3", Severity: "HIGH",
	}); got != want {
		t.Fatalf("vuln[0] = %+v, want %+v", got, want)
	}
	if result.Vulnerabilities[4].Severity != "UNKNOWN" {
		t.Fatalf("expected UNKNOWN severity, got %q", result.Vulnerabilities[4].Severity)
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
