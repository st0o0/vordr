# Vordr

Lightweight CVE scanner for homelab Docker stacks. Watches running containers via the Docker API, scans images via bundled Trivy CLI, and exposes vulnerability metrics via a Prometheus endpoint. Single container — no Trivy sidecar needed.

## Architecture

```
cmd/vordr/          Entry point, healthcheck, signal handling
internal/
  config/           VORDR_ env vars, testable via getenv function
  metrics/          Prometheus /metrics endpoint (vordr_ prefix), /healthz
  watcher/          Docker API container discovery + event stream
  scheduler/        Periodic rescan ticker
  worker/           N-goroutine scan worker pool reading from scan channel
  store/            Thread-safe image state (sync.RWMutex + map)
  trivy/            CLI subprocess wrapper (exec.Command trivy image --format json)
```

## Build

```bash
go build -o vordr ./cmd/vordr
go test ./...
go vet ./...
golangci-lint run
```

## Conventions

- Go 1.26+, pure Go, no CGO
- Conventional commits: feat/fix/perf/docs/chore/refactor/test/ci/build/deps
- commitlint enforced, release-please for versioning
- golangci-lint v2, hadolint for Dockerfile
- All env vars prefixed with `VORDR_`
- Structured logging via log/slog (JSON or text)
- Metric names prefixed with `vordr_`
- Custom Prometheus registry (never DefaultRegisterer)
- Errors wrapped with `fmt.Errorf` and `%w`
- Config loaded via `config.Load(os.Getenv)` — testable by injecting getenv

## Deployment

Single container with bundled Trivy binary (~200MB). Trivy DB cached on a volume at VORDR_TRIVY_CACHE_DIR. Deployed via compose override in homelab-infra Alloy stack. Published as `ghcr.io/st0o0/vordr`.
