# vordr

[![Release](https://img.shields.io/github/v/release/st0o0/vordr?sort=semver&style=flat-square)](https://github.com/st0o0/vordr/releases)
[![GHCR](https://img.shields.io/badge/ghcr.io-st0o0%2Fvordr-2496ED?style=flat-square&logo=docker&logoColor=white)](https://github.com/st0o0/vordr/pkgs/container/vordr)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg?style=flat-square)](LICENSE)

Lightweight CVE scanner for homelab Docker stacks. Watches running containers via the Docker API, scans images with a bundled Trivy CLI, and exposes vulnerability metrics via a Prometheus endpoint.

```yaml
services:
  vordr:
    image: ghcr.io/st0o0/vordr:latest
    restart: unless-stopped
    environment:
      VORDR_SCAN_INTERVAL: 6h
      VORDR_SCAN_WORKERS: 3
    ports:
      - "9550:9550"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - trivy-cache:/var/cache/trivy

volumes:
  trivy-cache:
```

## Features

- **Container discovery** -- lists running containers on startup and watches Docker events for start/stop in real time
- **Bundled Trivy** -- single container image (~200 MB), no sidecar needed; vulnerability DB cached on a volume
- **Prometheus metrics** -- per-image vulnerability counts by severity, scan status, timing, container-to-image mapping, and Trivy DB age
- **`/details` endpoint** -- JSON API with per-image CVE lists including package, installed version, fix version, and severity
- **Image exclusion** -- glob patterns to skip known-noisy base images
- **Delta detection** -- on rescan, only newly appeared CVEs are logged; severity threshold controls what is alert-worthy
- **Retry backoff** -- failed scans are retried with exponential backoff (30 s to 15 min) instead of waiting for the next scheduled rescan
- **Staggered startup** -- discovered images are queued with a configurable gap to avoid overwhelming Trivy on cold start
- **Grafana dashboard** -- importable JSON dashboard with stat panels, severity charts, per-image table, and DB age gauge

## Configuration

All via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `VORDR_TRIVY_CACHE_DIR` | `/var/cache/trivy` | Trivy vulnerability DB cache directory |
| `VORDR_TRIVY_TIMEOUT` | `5m` | Timeout per Trivy scan |
| `VORDR_SCAN_INTERVAL` | `6h` | Periodic rescan interval |
| `VORDR_SCAN_WORKERS` | `3` | Parallel scan workers |
| `VORDR_SCAN_STAGGER` | `30s` | Gap between scans at startup |
| `VORDR_LISTEN_ADDR` | `:9550` | HTTP listen address (metrics, health, details) |
| `VORDR_LOG_LEVEL` | `info` | Log level (`debug`, `info`, `warn`, `error`) |
| `VORDR_LOG_FORMAT` | `text` | Log format (`text`, `json`) |
| `VORDR_EXCLUDE_IMAGES` | _(empty)_ | Comma-separated glob patterns to skip (e.g. `aquasec/trivy:*,*devcontainer*`) |
| `VORDR_ALERT_SEVERITY` | `HIGH` | Minimum severity for delta-detection alerts |
| `VORDR_RETRY_MAX_BACKOFF` | `15m` | Maximum retry backoff for failed scans |

## How it works

```
Docker API ──▶ Watcher ──┐
                          ├──▶ scanCh ──▶ WorkerPool ──▶ Store
             Scheduler ──┘       ▲              │          │
                                 │              ▼          ▼
             RetryQueue ─────────┘         Trivy CLI   /metrics
                (backoff)                              /details
                                                       /healthz
```

1. **Watcher** lists all running containers on startup and subscribes to Docker events. New images are queued for scanning with staggered timing.
2. **WorkerPool** (N goroutines) reads from the scan channel, calls `trivy image --format json`, and updates the Store with results.
3. **Scheduler** re-queues all known images at the configured interval so newly disclosed CVEs are caught.
4. **RetryQueue** re-queues failed scans with exponential backoff.
5. **Store** holds thread-safe image state (sync.RWMutex) and serves snapshots to the metrics collector and /details handler.
6. **Metrics collector** generates Prometheus metrics on each scrape from the Store snapshot.

## Endpoints

| Path | Description |
|------|-------------|
| `GET /metrics` | Prometheus metrics (text format) |
| `GET /healthz` | JSON health check (`{"status":"ok"}`) |
| `GET /details` | JSON array of all images with CVE details |
| `GET /details?image=nginx:latest` | Filtered to a single image |

## Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `vordr_image_scan_status` | gauge | `image`, `tag` | Scan state (0=pending, 1=scanning, 2=scanned, 3=error) |
| `vordr_vulnerability_count` | gauge | `image`, `tag`, `severity` | Vulnerability count by severity |
| `vordr_fixable_count` | gauge | `image`, `tag`, `severity` | Fixable vulnerability count by severity |
| `vordr_scan_timestamp_seconds` | gauge | `image`, `tag` | Unix timestamp of last scan |
| `vordr_scan_duration_seconds` | gauge | `image`, `tag` | Duration of last scan in seconds |
| `vordr_container_image_info` | gauge | `container`, `image`, `tag` | Container-to-image mapping (value 1) |
| `vordr_images_discovered` | gauge | | Total tracked images |
| `vordr_images_scanned` | gauge | | Successfully scanned images |
| `vordr_images_pending` | gauge | | Images awaiting scan |
| `vordr_images_errors` | gauge | | Images in error state |
| `vordr_scan_workers_active` | gauge | | Currently active scan workers |
| `vordr_scan_queue_length` | gauge | | Current scan queue length |
| `vordr_trivy_db_updated_timestamp_seconds` | gauge | | Unix timestamp of last Trivy DB update |

## Image

- **Registry:** `ghcr.io/st0o0/vordr`
- **Tags:** semantic version (`1.0.0`), `latest`
- **Architectures:** `linux/amd64`

## Grafana dashboard

Import `grafana/vordr-dashboard.json` into Grafana. It includes stat panels, severity bar charts, a per-image CVE table, Trivy DB age gauge, and scan duration history. Select your Prometheus datasource on import.

## Development

```bash
go build -o vordr ./cmd/vordr
go test ./...
go vet ./...
golangci-lint run
docker build -t vordr:ci .
```

Commits follow [Conventional Commits](https://www.conventionalcommits.org/); releases and the GHCR image are cut automatically by release-please. See [CONTRIBUTING.md](CONTRIBUTING.md) for full development guidelines.

## License

MIT (see [LICENSE](LICENSE)). vordr is a pure-Go static binary; the built image bundles the Trivy scanner binary (Apache-2.0). Its Go module dependencies are permissively licensed (MIT/Apache-2.0/BSD); see [NOTICE](NOTICE).
