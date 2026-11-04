# server-image Specification

## Purpose
TBD - created by archiving change add-trivy-server-mode. Update Purpose after archive.
## Requirements
### Requirement: Companion server Dockerfile
The project SHALL provide a `Dockerfile.server` that builds a pre-configured Trivy server image.

#### Scenario: Server image build
- **WHEN** `docker build -f Dockerfile.server -t vordr-server .` is run
- **THEN** the resulting image runs Trivy in server mode listening on port 4954

#### Scenario: Server healthcheck
- **WHEN** the vordr-server container is running
- **THEN** the Docker healthcheck reports healthy via Trivy's built-in `/healthz` endpoint

#### Scenario: Version pin matches main image
- **WHEN** `Dockerfile.server` is inspected
- **THEN** the Trivy version matches the version used in the main `Dockerfile` trivy stage

### Requirement: Server image published to GHCR
The vordr-server image SHALL be published as `ghcr.io/st0o0/vordr-server` alongside the main vordr image.

#### Scenario: Release publishes both images
- **WHEN** a release is created via release-please
- **THEN** both `ghcr.io/st0o0/vordr` and `ghcr.io/st0o0/vordr-server` are built and pushed with the same version tag

