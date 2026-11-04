# e2e-dev-stack Specification

## Purpose
TBD - created by archiving change add-e2e-dev-testing. Update Purpose after archive.
## Requirements
### Requirement: Dev compose file defines vordr and test workloads
The project SHALL provide a `docker-compose.dev.yml` at the repository root that defines the vordr service built from source and three test workload containers.

#### Scenario: vordr service configuration
- **WHEN** `docker compose -f docker-compose.dev.yml up -d --build` is run
- **THEN** vordr is built from the repository Dockerfile, exposes port 9550, mounts `/var/run/docker.sock` read-only, and uses a named `trivy-cache` volume

#### Scenario: Dev environment variables
- **WHEN** the vordr service starts from docker-compose.dev.yml
- **THEN** the following env vars are set: `VORDR_SCAN_INTERVAL=30s`, `VORDR_SCAN_STAGGER=2s`, `VORDR_SCAN_WORKERS=2`, `VORDR_LOG_LEVEL=debug`, `VORDR_LOG_FORMAT=text`, `VORDR_EXCLUDE_IMAGES=redis*,vordr-*`

#### Scenario: Clean workload container
- **WHEN** the dev stack is running
- **THEN** a container named `workload-clean` runs `alpine:3.21` with a long-lived command (e.g., `sleep infinity`)

#### Scenario: Vulnerable workload container
- **WHEN** the dev stack is running
- **THEN** a container named `workload-vuln` runs `nginx:1.24.0`

#### Scenario: Excluded workload container
- **WHEN** the dev stack is running
- **THEN** a container named `workload-exclude` runs `redis:7`, matching the `VORDR_EXCLUDE_IMAGES=redis*,vordr-*` pattern

### Requirement: Trivy cache persists across restarts
The trivy-cache named volume SHALL persist across `docker compose down` (without `-v`) to avoid repeated DB downloads during development.

#### Scenario: Cache survival on restart
- **WHEN** `docker compose -f docker-compose.dev.yml down` followed by `up -d` is run
- **THEN** Trivy does not re-download the vulnerability database

#### Scenario: Clean slate with volume removal
- **WHEN** `docker compose -f docker-compose.dev.yml down -v` is run
- **THEN** the trivy-cache volume is removed and the next startup triggers a fresh DB download

