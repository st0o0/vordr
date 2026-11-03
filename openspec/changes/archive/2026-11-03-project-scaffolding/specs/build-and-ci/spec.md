## ADDED Requirements

### Requirement: Multi-stage Dockerfile
The project SHALL have a multi-stage Dockerfile that builds a static Go binary (CGO_ENABLED=0) and copies it into a `scratch` image. The Dockerfile SHALL support multi-arch builds via BUILDPLATFORM/TARGETARCH ARGs.

#### Scenario: Docker build succeeds
- **WHEN** running `docker build -t vordr .`
- **THEN** the build SHALL produce a working image under 15MB

#### Scenario: Version injection
- **WHEN** building with `--build-arg VERSION=1.0.0`
- **THEN** the resulting binary SHALL report version 1.0.0

#### Scenario: OCI labels
- **WHEN** inspecting the built image
- **THEN** it SHALL have org.opencontainers.image.title, .description, .source, .documentation, and .licenses labels

#### Scenario: Built-in healthcheck
- **WHEN** the container is running
- **THEN** Docker SHALL execute `/vordr healthcheck` every 30s with 10s timeout

### Requirement: CI workflow via reusable workflows
All CI workflows SHALL delegate to `st0o0/github-workflows` reusable workflows.

#### Scenario: CI on push/PR
- **WHEN** code is pushed or a PR is opened
- **THEN** `go-ci.yml` SHALL run (lint, test, build) and `commitlint.yml` SHALL validate commit messages

#### Scenario: Security scanning
- **WHEN** Dockerfile, go.mod, or go.sum changes (or weekly cron)
- **THEN** Trivy scan SHALL run against the repository

#### Scenario: Release automation
- **WHEN** a release is created by release-please
- **THEN** `docker-release.yml` SHALL build and push the image to GHCR

#### Scenario: Dev build
- **WHEN** a PR has the `dev-build` label
- **THEN** a dev image SHALL be built and pushed to GHCR

### Requirement: Linting configuration
The project SHALL use golangci-lint v2 with exclusion presets (comments, common-false-positives, legacy, std-error-handling) and hadolint with failure-threshold warning.

#### Scenario: Lint passes on scaffolded code
- **WHEN** running `golangci-lint run`
- **THEN** the scaffolded code SHALL produce zero lint errors

### Requirement: Release automation
The project SHALL use release-please with `simple` release type and conventional commit-based changelog generation.

#### Scenario: Initial manifest
- **WHEN** inspecting .release-please-manifest.json
- **THEN** it SHALL contain `{".":"0.0.0"}`

### Requirement: Dependency automation
The project SHALL use Renovate extending `local>st0o0/renovate-config:go` and `local>st0o0/renovate-config:docker`.

#### Scenario: Renovate config valid
- **WHEN** Renovate processes the repository
- **THEN** it SHALL detect Go and Docker dependencies and create update PRs
