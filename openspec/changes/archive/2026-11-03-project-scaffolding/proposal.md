## Why

Vordr hat eine detaillierte ROADMAP.md, aber noch kein Go-Projekt. Bevor die v0.1 MVP-Implementierung beginnen kann, braucht das Repo die gleiche Infrastruktur wie alle anderen Go-Repos (bifrost, eir, forseti, mjolnir, ran, var): CI/CD, Tooling-Configs, Dockerfile, Projektstruktur, Community-Files. Durch konsistentes Scaffolding sind CI, Releases und Dependency-Updates vom ersten Commit an funktionsfähig.

## What Changes

- **Go-Modul initialisieren**: `github.com/st0o0/vordr`, Go 1.26.x, Dependencies (docker/docker, prometheus/client_golang, golang.org/x/sync)
- **cmd/vordr/ Einstiegspunkt**: main.go (version, healthcheck subcommand, config.Load, slog, signal, prometheus.NewRegistry) + healthcheck.go — nach bewährtem Template
- **internal/ Package-Struktur**: Leere Packages mit Grundstruktur für config/, metrics/, watcher/, scheduler/, worker/, store/, trivy/
- **Dockerfile**: Multi-stage Build → scratch, CGO_ENABLED=0, ldflags Version-Injection, OCI Labels, Healthcheck
- **docker-compose.yml**: vordr + trivy server Sidecar
- **CI/CD Workflows**: ci.yml, release.yml, security.yml, dev-build.yml, labeler.yml, labels.yml — alle via st0o0/github-workflows
- **Tooling-Configs**: .golangci.yml, .hadolint.yaml, commitlint.config.mjs, release-please-config.json, .release-please-manifest.json, renovate.json
- **Repo-Files**: .gitignore, .dockerignore, LICENSE (MIT), NOTICE
- **Community-Files**: CODE_OF_CONDUCT.md, CONTRIBUTING.md, SECURITY.md
- **AI-Config**: CLAUDE.md (Projekt-Kontext, Build-Commands, Konventionen), AGENTS.md (Code-Change-Instruktionen, OpenSpec-Nutzung)
- **GitHub Issue-Templates**: bug_report.yml, feature_request.yml, config.yml
- **tests/ Verzeichnis**: Grundstruktur für Integration-Tests

## Capabilities

### New Capabilities
- `go-project-structure`: Go-Modul, cmd/ Einstiegspunkt mit main.go/healthcheck.go, internal/ Package-Layout
- `build-and-ci`: Dockerfile, CI/CD Workflows, Tooling-Configs (lint, hadolint, commitlint), Release-Automation
- `metrics-server`: Prometheus Custom Registry, HTTP-Server mit /metrics und /healthz auf :9550
- `config-loading`: Env-Var-basierte Konfiguration mit Load(getenv)-Pattern, Validierung, Defaults
- `deployment`: docker-compose.yml für vordr + trivy server Sidecar
- `repo-standards`: Community-Files, Issue-Templates, AI-Config (CLAUDE.md, AGENTS.md), LICENSE, NOTICE

### Modified Capabilities

_Keine — das Repo ist neu, es gibt keine bestehenden Specs._

## Impact

- **Neues Go-Modul** mit Dependencies auf docker/docker, prometheus/client_golang, golang.org/x/sync
- **GitHub Actions** Workflows referenzieren st0o0/github-workflows (muss als Repo existieren)
- **GHCR** Docker-Images werden via CI gebaut und gepusht
- **Renovate** überwacht Go- und Docker-Dependencies automatisch
- **Release-Please** managed Versioning und Changelog
- Die internal/ Packages sind zunächst Stubs — die eigentliche Implementierung folgt als separater Change
