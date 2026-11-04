# e2e-verify-skill Specification

## Purpose
TBD - created by archiving change add-e2e-dev-testing. Update Purpose after archive.
## Requirements
### Requirement: Skill orchestrates full test plan execution
The `.claude/skills/e2e-verify/skill.md` skill SHALL read `E2E-TEST-PLAN.md` and execute all test sections sequentially, producing a results file.

#### Scenario: Skill invocation
- **WHEN** the user runs `/e2e-verify`
- **THEN** the skill reads the full test plan, manages the dev stack lifecycle, executes all tests, and writes results

### Requirement: Skill manages dev stack lifecycle
The skill SHALL handle compose up/down as part of test execution.

#### Scenario: Clean start
- **WHEN** the skill starts
- **THEN** it runs `docker compose -f docker-compose.dev.yml down -v` followed by `up -d --build` for a clean slate

#### Scenario: Health poll before tests
- **WHEN** the stack is started
- **THEN** the skill polls `GET http://localhost:9550/healthz` until it returns 200, with a timeout of 5 minutes to accommodate Trivy DB download

### Requirement: Skill waits for initial scans to complete
The skill SHALL not begin functional tests until vordr has completed scanning the test workloads.

#### Scenario: Scan completion detection
- **WHEN** healthz returns 200
- **THEN** the skill polls `GET http://localhost:9550/metrics` until metrics for both alpine and nginx images appear with scan results, with a 5-minute timeout

### Requirement: All tests produce PASS or FAIL
Every test in the plan SHALL have a definitive outcome. SKIP is not a valid status.

#### Scenario: No skipped tests
- **WHEN** the results file is written
- **THEN** every test has status PASS or FAIL, never SKIP

#### Scenario: Setup failure cascades
- **WHEN** a setup step fails (e.g., compose up fails)
- **THEN** that step and all dependent tests are marked FAIL with an explanatory note

### Requirement: Results file format
The results file SHALL follow a consistent structure with a summary and per-section details.

#### Scenario: Results file content
- **WHEN** test execution completes
- **THEN** `E2E-TEST-RESULTS-<DATE>.md` contains a summary table (total PASS/FAIL counts) and a per-section table with columns: #, Test, Status, Notes

### Requirement: Tests use curl or Invoke-RestMethod
All test verification SHALL use HTTP requests to vordr endpoints and docker compose commands. No browser automation is needed.

#### Scenario: API-only testing
- **WHEN** tests verify vordr behavior
- **THEN** they use `curl` or `Invoke-RestMethod` for HTTP endpoints and `docker compose` commands for container lifecycle operations

