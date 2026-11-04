## Why

vordr has 12 unit test files covering individual components with hand-rolled mocks, but no integration or end-to-end tests. The critical interaction paths - Docker API event handling, Trivy CLI subprocess scanning, container lifecycle tracking through metrics - are untested. A local dev stack with a reproducible test plan lets us verify vordr's real behavior before shipping changes.

## What Changes

- Add `docker-compose.dev.yml` with vordr (built from source) and three test workload containers (clean, vulnerable, excluded)
- Add `E2E-TEST-PLAN.md` as a reusable test specification covering ~8 test areas: initial scan, container lifecycle, exclude logic, rescan cycle, error handling, healthcheck, metrics shape, and details endpoint
- Add `.claude/skills/e2e-verify/skill.md` as an orchestrator skill that executes the test plan, producing `E2E-TEST-RESULTS-<DATE>.md` output files

## Capabilities

### New Capabilities
- `e2e-dev-stack`: Docker Compose dev environment with vordr and test workloads for local E2E verification
- `e2e-test-plan`: Reusable numbered test specification covering all observable vordr behaviors
- `e2e-verify-skill`: Claude Code skill that orchestrates test plan execution and produces structured results

### Modified Capabilities

None - this is additive testing infrastructure with no changes to existing specs.

## Impact

- No changes to application code
- New files: `docker-compose.dev.yml` (repo root), `E2E-TEST-PLAN.md` (repo root), `.claude/skills/e2e-verify/skill.md`
- Generated files: `E2E-TEST-RESULTS-*.md` (gitignored)
- Dependencies: requires Docker Desktop/Engine with compose, existing Dockerfile for build
- Timing consideration: first run downloads Trivy vulnerability DB (~40MB, 1-2 min)
