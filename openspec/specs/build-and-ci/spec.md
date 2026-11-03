## MODIFIED Requirements

### Requirement: Multi-stage Dockerfile
The project SHALL have a multi-stage Dockerfile that: (1) builds the vordr Go binary, (2) copies the Trivy binary from the official aquasec/trivy image, and (3) produces a final Alpine-based image with both binaries, ca-certificates, and the healthcheck.

#### Scenario: Docker build succeeds
- **WHEN** running `docker build -t vordr .`
- **THEN** the build SHALL produce a working image under 250MB

#### Scenario: Trivy binary available
- **WHEN** the container starts
- **THEN** `trivy --version` SHALL succeed inside the container

#### Scenario: Version injection
- **WHEN** building with `--build-arg VERSION=1.0.0`
- **THEN** the resulting vordr binary SHALL report version 1.0.0

#### Scenario: OCI labels
- **WHEN** inspecting the built image
- **THEN** it SHALL have org.opencontainers.image.title, .description, .source, .documentation, and .licenses labels

#### Scenario: Built-in healthcheck
- **WHEN** the container is running
- **THEN** Docker SHALL execute `/vordr healthcheck` every 30s with 10s timeout

#### Scenario: Pinned Trivy version
- **WHEN** inspecting the Dockerfile
- **THEN** the Trivy image SHALL be pinned to a specific version tag (not `latest`)
