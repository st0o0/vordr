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
	cacheDir string
	timeout  time.Duration
}

func NewCLI(cacheDir string, timeout time.Duration) *CLI {
	return &CLI{
		cacheDir: cacheDir,
		timeout:  timeout,
	}
}

type ScanResult struct {
	Counts  map[string]int // severity → total count
	Fixable map[string]int // severity → fixable count
}

type trivyReport struct {
	Results []trivyResult `json:"Results"`
}

type trivyResult struct {
	Vulnerabilities []vulnerability `json:"Vulnerabilities"`
}

type vulnerability struct {
	Severity     string `json:"Severity"`
	FixedVersion string `json:"FixedVersion"`
}

func (c *CLI) Scan(ctx context.Context, imageRef string) (*ScanResult, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "trivy", "image",
		"--format", "json",
		"--scanners", "vuln",
		"--quiet",
		"--cache-dir", c.cacheDir,
		imageRef)

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
		}
	}

	return result, nil
}
