## 1. Go Module & Project Structure

- [x] 1.1 Initialize go.mod (`github.com/st0o0/vordr`, Go 1.26.x)
- [x] 1.2 Add dependencies: docker/docker, prometheus/client_golang, golang.org/x/sync
- [x] 1.3 Create `cmd/vordr/main.go` with version var, subcommand dispatch, config.Load, slog logger, signal handling, errgroup, metrics server startup
- [x] 1.4 Create `cmd/vordr/healthcheck.go` with HTTP GET to /healthz and exit code logic
- [x] 1.5 Create stub packages: internal/watcher/, internal/scheduler/, internal/worker/, internal/store/, internal/trivy/ (package declaration only)
- [x] 1.6 Verify `go build ./...` compiles without errors

## 2. Config Package

- [x] 2.1 Create `internal/config/config.go` with Config struct and Load(getenv) function
- [x] 2.2 Implement env var parsing: VORDR_TRIVY_ADDR, VORDR_SCAN_INTERVAL, VORDR_SCAN_WORKERS, VORDR_SCAN_STAGGER, VORDR_LISTEN_ADDR, VORDR_LOG_LEVEL, VORDR_LOG_FORMAT
- [x] 2.3 Add defaults, duration/int parsing helpers, validation with descriptive errors

## 3. Metrics Package

- [x] 3.1 Create `internal/metrics/metrics.go` with New(reg) constructor and custom Prometheus registry
- [x] 3.2 Create `internal/metrics/server.go` with HTTP server (/metrics + /healthz), ReadHeaderTimeout 5s, graceful shutdown

## 4. Dockerfile & Compose

- [x] 4.1 Create multi-stage Dockerfile: golang:alpine build → scratch, CGO_ENABLED=0, ldflags version injection, OCI labels, healthcheck CMD
- [x] 4.2 Create docker-compose.yml with vordr + trivy server, Docker socket mount (read-only), trivy-cache volume

## 5. Tooling Configs (1:1 from template)

- [x] 5.1 Create .golangci.yml (v2, exclusion presets)
- [x] 5.2 Create .hadolint.yaml (failure-threshold warning, ignore DL3018)
- [x] 5.3 Create commitlint.config.mjs (conventional commits + deps type)
- [x] 5.4 Create release-please-config.json + .release-please-manifest.json (version 0.0.0)
- [x] 5.5 Create renovate.json (extends go + docker configs)
- [x] 5.6 Create .gitignore (vordr binary, coverage, IDE, .env, .secrets)
- [x] 5.7 Create .dockerignore (.git, .github, tests, docs, *.md except LICENSE, tooling configs)

## 6. CI/CD Workflows

- [x] 6.1 Create .github/workflows/ci.yml (go-ci + commitlint, docker-image-name: vordr)
- [x] 6.2 Create .github/workflows/release.yml (release-please + conditional docker-release)
- [x] 6.3 Create .github/workflows/security.yml (Trivy scan on Dockerfile/go.mod/go.sum + weekly cron)
- [x] 6.4 Create .github/workflows/dev-build.yml (label-gated dev builds to GHCR)
- [x] 6.5 Create .github/workflows/labeler.yml + labels.yml (PR auto-labeling + weekly sync)

## 7. Community & Repo Files

- [x] 7.1 Create LICENSE (MIT)
- [x] 7.2 Create NOTICE (copyright + key dependencies)
- [x] 7.3 Create CODE_OF_CONDUCT.md (Contributor Covenant v2.1)
- [x] 7.4 Create CONTRIBUTING.md (dev setup, build/test/lint commands, PR rules with conventional commits)
- [x] 7.5 Create SECURITY.md (private reporting via GitHub advisories)

## 8. GitHub Issue Templates

- [x] 8.1 Create .github/ISSUE_TEMPLATE/bug_report.yml (what happened, repro, logs, image tag, host info)
- [x] 8.2 Create .github/ISSUE_TEMPLATE/feature_request.yml (problem, proposal, alternatives)
- [x] 8.3 Create .github/ISSUE_TEMPLATE/config.yml (blank issues disabled, security → advisories)

## 9. AI Configuration

- [x] 9.1 Create CLAUDE.md (project purpose, architecture overview, package layout, build/test/lint commands, conventions, deployment context)
- [x] 9.2 Create AGENTS.md (read CLAUDE.md first, run build+test+lint after changes, OpenSpec usage, security-sensitive areas)

## 10. Integration Test Structure

- [x] 10.1 Create tests/ directory with a placeholder test file
- [x] 10.2 Verify full build pipeline: `go build`, `go vet`, `golangci-lint run` all pass
