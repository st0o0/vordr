---
description: "Run the full E2E test plan against the dev stack. Use when the user wants to verify all features, run E2E tests, or do a final verification pass."
---

# E2E Verify

Execute `E2E-TEST-PLAN.md` against the dev Docker stack. Every test must be PASS or FAIL, no SKIP.

## How to Run

1. Read `E2E-TEST-PLAN.md` fully before starting
2. Execute Phase 0 (stack lifecycle)
3. Follow the test plan sections 0-8 in order
4. Write results to `E2E-TEST-RESULTS-<DATE>.md`

## Phase 0: Stack Lifecycle

### Reset and Start

```powershell
docker compose -f docker-compose.dev.yml down -v
docker compose -f docker-compose.dev.yml up -d --build
```

### Wait for Health

Poll `GET http://localhost:9550/healthz` every 5 seconds until it returns HTTP 200. Timeout: 3 minutes.

If healthz does not respond within 3 minutes, check logs:
```powershell
docker compose -f docker-compose.dev.yml logs vordr --tail 100
```
Mark Section 0 as FAIL and stop.

### Wait for Initial Scans

After healthz responds, poll `GET http://localhost:9550/metrics` every 10 seconds until `vordr_images_scanned` shows value 2 (both alpine and nginx scanned). Timeout: 5 minutes.

This accounts for Trivy DB download on first run (~40MB, 1-2 min) plus scan time.

If scans do not complete within 5 minutes, check logs and mark Section 1 as FAIL.

## Test Execution

Execute each section from `E2E-TEST-PLAN.md` in order (0 through 8). Use `Invoke-RestMethod` or `curl` for HTTP requests and `docker compose` commands for container operations.

### Test Methods

**HTTP endpoints:**
```powershell
# Health check
Invoke-RestMethod http://localhost:9550/healthz

# Prometheus metrics (raw text)
(Invoke-WebRequest http://localhost:9550/metrics).Content

# Details endpoint (JSON)
Invoke-RestMethod http://localhost:9550/details
Invoke-RestMethod "http://localhost:9550/details?image=nginx:1.24.0"
```

**Container operations:**
```powershell
docker compose -f docker-compose.dev.yml stop workload-vuln
docker compose -f docker-compose.dev.yml start workload-vuln
docker compose -f docker-compose.dev.yml ps --format json
docker inspect vordr --format "{{.State.Health.Status}}"
```

**Logs:**
```powershell
docker compose -f docker-compose.dev.yml logs vordr --tail 50
```

### Polling Pattern

For tests that require waiting (scan completion, rescan cycle):
1. Record the initial state
2. Wait the expected duration (+ 5s margin)
3. Check the new state
4. If not met, retry up to 3 times with 10s intervals
5. If still not met, mark as FAIL with the observed values

### Section Dependencies

- Section 0 must PASS before continuing to any other section
- Section 1 requires scans to complete (poll metrics)
- Section 2 modifies container state (stop/start workload-vuln) - must restore state before Section 3
- Section 4 requires a ~35 second wait for rescan interval
- All other sections can run on existing state

### Error Cascade

If a setup step fails (compose up, healthz timeout), mark that step and all dependent tests as FAIL with a note explaining the dependency.

## Results File

Write `E2E-TEST-RESULTS-YYYY-MM-DD.md` in the repository root.

### Format

```markdown
# E2E Test Results

**Date:** YYYY-MM-DD
**Stack:** docker-compose.dev.yml (build from source)
**Duration:** Xm Ys

## Summary

| Status | Count |
|--------|-------|
| PASS   | N     |
| FAIL   | N     |
| **Total** | **N** |

## Results

### Section 0: Setup Verification

| # | Test | Status | Notes |
|---|------|--------|-------|
| 0.1 | All containers running | PASS | 4/4 containers up |
| 0.2 | Health endpoint responds | PASS | HTTP 200, version: dev |
| 0.3 | Logs indicate startup | PASS | |

### Section 1: Initial Scan

| # | Test | Status | Notes |
|---|------|--------|-------|
| 1.1 | Non-excluded images in metrics | PASS | discovered=2, scanned=2 |
| ... | ... | ... | ... |
```

Continue this pattern for all sections.

### Notes Column

- PASS: brief confirmation or key values observed
- FAIL: what was expected vs what was observed

## Rules

- Do NOT modify `E2E-TEST-PLAN.md` - it is the reusable spec
- Write results to a separate file only
- Follow the execution order exactly
- If setup fails, mark dependent tests as FAIL
- Every test must be PASS or FAIL, never SKIP
- Include actual observed values in the Notes column for key assertions
