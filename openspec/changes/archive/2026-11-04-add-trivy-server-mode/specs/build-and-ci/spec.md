## ADDED Requirements

### Requirement: CI builds and publishes Docker images
The CI pipeline SHALL build and publish both `ghcr.io/st0o0/vordr` and `ghcr.io/st0o0/vordr-server` images.

#### Scenario: Release workflow publishes both images
- **WHEN** release-please creates a new release
- **THEN** the release workflow builds and pushes `ghcr.io/st0o0/vordr` from `Dockerfile`
- **AND** builds and pushes `ghcr.io/st0o0/vordr-server` from `Dockerfile.server`
- **AND** both images use the same version tag

#### Scenario: Dev-build workflow builds both images
- **WHEN** a PR has the `dev-build` label
- **THEN** the dev-build workflow builds and pushes dev images for both vordr and vordr-server

#### Scenario: Trivy version consistency
- **WHEN** both Dockerfiles are inspected
- **THEN** they reference the same `aquasec/trivy` version tag
