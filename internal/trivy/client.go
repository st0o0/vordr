package trivy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

type CLI struct {
	cacheDir  string
	serverURL string
	timeout   time.Duration
}

func NewCLI(cacheDir string, timeout time.Duration, serverURL string) *CLI {
	return &CLI{
		cacheDir:  cacheDir,
		serverURL: serverURL,
		timeout:   timeout,
	}
}

type ScanResult struct {
	Counts          map[string]int // severity → total count
	Fixable         map[string]int // severity → fixable count
	Vulnerabilities []Vulnerability
}

type Vulnerability struct {
	ID               string
	PkgName          string
	InstalledVersion string
	FixedVersion     string
	Severity         string
}

type trivyReport struct {
	Results []trivyResult `json:"Results"`
}

type trivyResult struct {
	Vulnerabilities []vulnerability `json:"Vulnerabilities"`
}

type vulnerability struct {
	VulnerabilityID  string `json:"VulnerabilityID"`
	PkgName          string `json:"PkgName"`
	InstalledVersion string `json:"InstalledVersion"`
	FixedVersion     string `json:"FixedVersion"`
	Severity         string `json:"Severity"`
}

func (c *CLI) Scan(ctx context.Context, imageRef string) (*ScanResult, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	args := []string{"image", "--format", "json", "--scanners", "vuln", "--quiet"}
	if c.serverURL != "" {
		args = append(args, "--server", c.serverURL)
	} else {
		args = append(args, "--cache-dir", c.cacheDir, "--skip-java-db-update")
	}
	args = append(args, imageRef)

	cmd := exec.CommandContext(ctx, "trivy", args...)

	out, err := cmd.Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("scan timed out after %s", c.timeout)
		}
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil, fmt.Errorf("trivy exited %d: %s", exitErr.ExitCode(), string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("trivy exec: %w", err)
	}

	var report trivyReport
	if err := json.Unmarshal(out, &report); err != nil {
		return nil, fmt.Errorf("parse trivy output: %w", err)
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

	return result, nil
}
