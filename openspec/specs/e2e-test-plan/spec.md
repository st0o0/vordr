# e2e-test-plan Specification

## Purpose
TBD - created by archiving change add-e2e-dev-testing. Update Purpose after archive.
## Requirements
### Requirement: Test plan document covers all observable behaviors
The project SHALL provide an `E2E-TEST-PLAN.md` at the repository root with numbered test sections covering vordr's complete observable behavior.

#### Scenario: Plan exists and is self-contained
- **WHEN** a user reads `E2E-TEST-PLAN.md`
- **THEN** it contains prerequisites, setup instructions, numbered test sections, and expected outcomes sufficient to execute the tests without additional documentation

### Requirement: Setup verification tests
Section 0 SHALL verify the dev stack is operational before running functional tests.

#### Scenario: Health check after startup
- **WHEN** the dev stack starts
- **THEN** `GET /healthz` returns HTTP 200

#### Scenario: Log output verification
- **WHEN** vordr starts
- **THEN** `docker compose logs vordr` contains structured log entries indicating startup

### Requirement: Initial scan tests
Section 1 SHALL verify vordr discovers and scans running containers on startup.

#### Scenario: Non-excluded images appear in metrics
- **WHEN** initial scans complete
- **THEN** `GET /metrics` contains `vordr_` metrics for both `alpine` and `nginx` images

#### Scenario: Vulnerable image has CVE counts
- **WHEN** nginx:1.24.0 is scanned
- **THEN** metrics show vulnerability counts greater than zero for at least one severity level

#### Scenario: Clean image has zero or minimal CVEs
- **WHEN** alpine:3.21 is scanned
- **THEN** metrics show zero or very low vulnerability counts

### Requirement: Container lifecycle tests
Section 2 SHALL verify vordr reacts to container stop and start events.

#### Scenario: Container stop removes image from tracking
- **WHEN** `docker compose -f docker-compose.dev.yml stop workload-vuln` is executed
- **THEN** subsequent `GET /metrics` no longer contains metrics for the nginx image

#### Scenario: Container start triggers rescan
- **WHEN** `docker compose -f docker-compose.dev.yml start workload-vuln` is executed
- **THEN** the nginx image reappears in `/metrics` with scan results after the scan completes

### Requirement: Exclude logic tests
Section 3 SHALL verify excluded images are never scanned or tracked.

#### Scenario: Excluded image absent from metrics
- **WHEN** the dev stack is running with `VORDR_EXCLUDE_IMAGES=redis*`
- **THEN** `GET /metrics` never contains metrics for the redis image at any point

### Requirement: Rescan cycle tests
Section 4 SHALL verify periodic rescanning works.

#### Scenario: Scan timestamp updates after interval
- **WHEN** at least `VORDR_SCAN_INTERVAL` (30s) has elapsed since the last scan
- **THEN** the scan timestamp in metrics or details endpoint reflects a more recent scan

### Requirement: Error handling tests
Section 5 SHALL verify vordr handles scan failures gracefully.

#### Scenario: Invalid image reference produces error state
- **WHEN** a container with an unresolvable image reference is started
- **THEN** the image appears in the store with an error status and the error is visible in details or logs

### Requirement: Healthcheck endpoint tests
Section 6 SHALL verify the healthcheck mechanism.

#### Scenario: HTTP healthcheck returns 200
- **WHEN** `GET /healthz` is called while vordr is running
- **THEN** the response status is 200

#### Scenario: Docker healthcheck reports healthy
- **WHEN** `docker inspect` is run on the vordr container
- **THEN** the health status is "healthy"

### Requirement: Metrics shape tests
Section 7 SHALL verify the Prometheus metrics have correct names and labels.

#### Scenario: All vordr metrics present
- **WHEN** `GET /metrics` is called after scans complete
- **THEN** the response contains all expected `vordr_` metric families

#### Scenario: Metric labels are correct
- **WHEN** vulnerability metrics are emitted for a scanned image
- **THEN** labels include `image` and `tag` with correct values, and severity-based metrics include the `severity` label

### Requirement: Details endpoint tests
Section 8 SHALL verify the vulnerability details HTTP endpoint.

#### Scenario: Details returns vulnerability data
- **WHEN** `GET /details` is called after scans complete
- **THEN** the response includes vulnerability details for scanned images

#### Scenario: Details filters by image
- **WHEN** `GET /details?image=<ref>` is called with a specific image reference
- **THEN** only vulnerability details for that image are returned

### Requirement: Results are written to separate file
Test results SHALL be written to `E2E-TEST-RESULTS-<DATE>.md`, never modifying the test plan itself.

#### Scenario: Results file naming
- **WHEN** the test plan is executed
- **THEN** results are written to `E2E-TEST-RESULTS-YYYY-MM-DD.md` in the repository root

#### Scenario: Plan remains unmodified
- **WHEN** tests are executed
- **THEN** `E2E-TEST-PLAN.md` has no uncommitted changes

