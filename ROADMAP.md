# Vördr -- CVE Scanner for Homelab Docker Stacks

Lightweight CVE scanner that watches running Docker containers and exposes vulnerability metrics via a Prometheus endpoint. Built with idiomatic Go concurrency (errgroup, channels, context) -- no actor framework. Designed to integrate into an existing Grafana/VictoriaMetrics/Alloy observability stack.

**Name:** Vördr (Old Norse) -- a guardian spirit that silently watches over its ward.

## Architecture

Vördr runs as a sidecar to a Trivy server instance. Vördr handles discovery, scheduling, and metrics. Trivy handles the actual vulnerability scanning.

```
┌─ compose.override.vordr.yaml ──────────────────────────────────────┐
│                                                                    │
│  ┌─ vordr (Go, distroless, ~10MB) ─────────────────────────────┐   │
│  │                                                             │   │
│  │  Docker API ──▶ discover running containers                 │   │
│  │  Docker Events ──▶ react to container start/stop            │   │
│  │  errgroup ──▶ structured concurrency                        │   │
│  │  HTTP :8080 ──▶ /metrics, /health, /details                 │   │
│  │                                                             │   │
│  │  Scan requests via HTTP ────────────────────────────┐       │   │
│  └─────────────────────────────────────────────────────┼───────┘   │
│                                                        │           │
│  ┌─ trivy server (:4954) ──────────────────────────────┼───────┐   │
│  │                                                     ▼       │   │
│  │  Receives scan requests, returns JSON results               │   │
│  │  Auto-updates vulnerability DB                              │   │
│  │  Caches scan results internally                             │   │
│  │  No Docker socket needed                                    │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                    │
│  volumes:                                                          │
│    trivy-cache: (persistent, ~50MB vulnerability DB)               │
│                                                                    │
└────────────────────────────────────────────────────────────────────┘
```

### Why two containers?

- **Independent updates** -- bump Trivy tag without rebuilding Vördr, and vice versa
- **Clean separation** -- Vördr is pure Go with no external binary dependencies
- **Trivy server caching** -- repeated scans of the same image are instant
- **Small Vördr image** -- ~10MB distroless vs ~180MB if bundling Trivy

## Concurrency Model

Idiomatic Go: errgroup for lifecycle, channels for communication, sync.RWMutex for shared state. No actor framework -- Go's built-in concurrency primitives map directly to the problem.

### .NET Analogy

```
.NET                              Go (Vördr)
──────────────────────            ──────────────────────
BackgroundService                 goroutine + context.Context
Channel<T>                        chan ImageRef
Task.WhenAll                      errgroup.Group
CancellationToken                 context.Context
SemaphoreSlim(N)                  N workers reading from channel
ConcurrentDictionary              sync.RWMutex + map
```

### Component Lifecycle

```go
func main() {
    ctx, cancel := signal.NotifyContext(context.Background(),
        syscall.SIGINT, syscall.SIGTERM)
    defer cancel()

    g, ctx := errgroup.WithContext(ctx)

    g.Go(func() error { return httpServer.Run(ctx) })
    g.Go(func() error { return watcher.Run(ctx) })
    g.Go(func() error { return scheduler.Run(ctx) })
    g.Go(func() error { return workerPool.Run(ctx) })

    if err := g.Wait(); err != nil {
        log.Fatal(err)
    }
}
```

All components share the same context. If one fails, the context is cancelled and everything shuts down gracefully. errgroup.Wait blocks until all goroutines return.

### Data Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│  Watcher                     Store                 HTTPServer   │
│  ┌──────────┐               ┌──────────────┐      ┌──────────┐ │
│  │          │               │ sync.RWMutex │      │          │ │
│  │ Docker   │──scanCh──▶   │              │ ◀────│ GET      │ │
│  │ events + │  chan ImageRef│ map[string]  │ .Snap│/metrics  │ │
│  │ initial  │               │  *ImageState │ shot │          │ │
│  │ listing  │               │              │ ()   │          │ │
│  └──────────┘               └──────┬───────┘      └──────────┘ │
│                                    │                            │
│  Scheduler                         │                            │
│  ┌──────────┐                      │                            │
│  │          │──scanCh──▶          │                            │
│  │ ticker   │  chan ImageRef       │                            │
│  │ 6h       │                      │                            │
│  └──────────┘                      │                            │
│                                    │                            │
│                             WorkerPool                          │
│                             ┌──────────────────────┐            │
│                             │                      │            │
│                             │ N goroutines reading  │            │
│                             │ from scanCh:          │            │
│                             │                      │            │
│                             │ for ref := range ch { │            │
│                             │   result := trivy(ref)│            │
│                             │   store.Update(result)│            │
│                             │ }                    │            │
│                             │                      │            │
│                             └──────────────────────┘            │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘

scanCh is a single channel shared between Watcher and Scheduler (producers)
and the WorkerPool (consumers). Channel capacity bounds the queue.
```

### Store

The store holds all image state and is the single source of truth. Workers write after each scan, the HTTP handler reads on every /metrics scrape.

```go
type Store struct {
    mu     sync.RWMutex
    images map[string]*ImageState
}

func (s *Store) Update(key string, state *ImageState) { /* Lock, write, Unlock */ }
func (s *Store) Remove(key string)                     { /* Lock, delete, Unlock */ }
func (s *Store) Snapshot() []ImageState                { /* RLock, copy, RUnlock */ }
```

Access pattern is simple: infrequent writes (after each scan, every few minutes), frequent reads (every Alloy scrape, every 30s). RWMutex is the right primitive -- no actor or channel indirection needed.

### Worker Pool

Bounded parallelism via N goroutines reading from a shared channel:

```go
func (p *Pool) Run(ctx context.Context) error {
    g, ctx := errgroup.WithContext(ctx)

    for i := 0; i < p.workers; i++ {
        g.Go(func() error {
            for {
                select {
                case <-ctx.Done():
                    return nil
                case ref, ok := <-p.scanCh:
                    if !ok {
                        return nil
                    }
                    p.scan(ctx, ref)  // HTTP to Trivy, update store
                }
            }
        })
    }

    return g.Wait()
}
```

Error handling per scan: log the error, set image state to Error, continue with next image. A single failed scan never crashes a worker.

## Scan State Model

Images are tracked at the **image level**, not container level. A separate info metric maps containers to images.

### States

```
Pending   -- discovered, waiting for scan worker
Scanning  -- scan in progress
Scanned   -- scan complete, results available
Error     -- scan failed, will retry on next cycle
```

### Clean startup behavior

On startup, Vördr lists all running containers and queues them for scanning with staggered scheduling (configurable gap between scans, default 30s). This avoids a load spike while filling in data progressively. Until an image has been scanned, no vulnerability metrics are emitted for it -- only the state metric shows it exists.

### Container tracking

```go
type ImageState struct {
    Status     Status            // Pending, Scanning, Scanned, Error
    Containers []string          // container names using this image
    Counts     map[string]int    // severity -> count (nil if not yet scanned)
    Fixable    map[string]int    // severity -> fixable count
    ScannedAt  time.Time         // zero if never scanned
    Duration   time.Duration     // last scan duration
    Error      string            // last error message
}
```

When a container starts with an already-known image: add to Containers list, no rescan needed. When the last container for an image stops: remove from store, metrics disappear.

## Prometheus Metrics

```prometheus
# Image scan state (0=pending, 1=scanning, 2=scanned, 3=error)
vordr_image_scan_status{image="grafana/grafana",tag="13.3.0"} 2

# Vulnerability counts (only emitted when state=scanned)
vordr_vulnerability_count{image="grafana/grafana",tag="13.3.0",severity="CRITICAL"} 0
vordr_vulnerability_count{image="grafana/grafana",tag="13.3.0",severity="HIGH"} 3
vordr_vulnerability_count{image="grafana/grafana",tag="13.3.0",severity="MEDIUM"} 12
vordr_vulnerability_count{image="grafana/grafana",tag="13.3.0",severity="LOW"} 27

# Fixable subset (fix version available in upstream)
vordr_fixable_count{image="grafana/grafana",tag="13.3.0",severity="HIGH"} 2

# Scan timing
vordr_scan_timestamp_seconds{image="grafana/grafana",tag="13.3.0"} 1730000000
vordr_scan_duration_seconds{image="grafana/grafana",tag="13.3.0"} 14.2

# Container-to-image mapping
vordr_container_image_info{container="grafana",image="grafana/grafana",tag="13.3.0"} 1

# Operational
vordr_images_discovered 42
vordr_images_pending 0
vordr_images_scanned 42
vordr_images_errors 0
vordr_scan_workers_active 0
vordr_scan_queue_length 0
vordr_trivy_db_updated_timestamp_seconds 1730000000
```

## Integration with Alloy Stack

Deployed as part of the Alloy stack via compose override:

```
stacks/alloy/
├── compose.override.vordr.yaml    # vordr + trivy server containers
```

Alloy scrapes the metrics endpoint via a feature config:

```
alloy/features/
├── vordr.alloy                    # Scrape config
```

No shared volumes needed -- pure HTTP scrape.

## Configuration

All via environment variables:

```
VORDR_TRIVY_CACHE_DIR=/var/cache/trivy  # Trivy DB cache directory
VORDR_TRIVY_TIMEOUT=5m                  # Timeout per scan
VORDR_SCAN_INTERVAL=6h                  # Periodic rescan interval
VORDR_SCAN_WORKERS=3                    # Parallel scan workers
VORDR_SCAN_STAGGER=30s                  # Gap between scans at startup
VORDR_LISTEN_ADDR=:9550                 # HTTP listen address
VORDR_LOG_LEVEL=info                    # Log level
```

## Design Decisions

- **Trivy CLI bundled, not as sidecar** -- single container, Trivy binary copied via multi-stage build. Research confirmed this is how Trivy Operator, Harbor Scanner, and GitLab do it. DB cached on a volume
- **Idiomatic Go, no actor framework** -- Go's channels, errgroup, and context.Context map directly to the concurrency needs. ProtoActor was considered but rejected: it's still beta, adds a large dependency, and the problem doesn't need distribution or clustering. The same patterns (worker pool, state machine, graceful shutdown) are achievable with ~200 LOC of standard Go.
- **Docker API for discovery, not compose-file parsing** -- we want to see what's actually running, not what's declared. A stopped container has no attack surface.
- **Event-driven + periodic hybrid** -- Docker events catch new deployments immediately, periodic rescans catch newly disclosed CVEs against the same images.
- **Image-level metrics with container mapping** -- deduplicates scan work (5 containers on same image = 1 scan), keeps cardinality bounded, Grafana can join via info metric.
- **No vulnerability metrics until scanned** -- avoids false "0 CVEs" for images that haven't been checked yet. Scan state metric makes the distinction explicit.
- **Prometheus endpoint, not textfile collector** -- cleaner integration, no shared volumes, Alloy can scrape directly.
- **Severity counts, not individual CVE labels** -- individual CVEs as Prometheus labels would create massive cardinality. Detail queries go through the /details JSON endpoint.

## Roadmap

### v0.1 -- MVP

- [ ] Go project scaffolding (module, Dockerfile, CI)
- [ ] Store: thread-safe image state management (sync.RWMutex)
- [ ] Watcher: Docker API container discovery + event stream
- [ ] Scheduler: periodic rescan ticker
- [ ] WorkerPool: N goroutines reading from scan channel, HTTP to Trivy
- [ ] Trivy client: HTTP client for Trivy server scan API + JSON parsing
- [ ] HTTP server: /metrics (Prometheus), /health
- [ ] Staggered startup scan (configurable gap)
- [ ] Graceful shutdown via context + errgroup
- [ ] Multi-stage Dockerfile (distroless)
- [ ] Compose override for homelab-infra Alloy stack (vordr + trivy server)
- [ ] Alloy scrape config

### v0.2 -- Observability

- [ ] /details JSON endpoint -- per-image CVE list with package, installed version, fix version
- [ ] Trivy DB age metric + alert when DB is stale
- [ ] Grafana dashboard (CVE overview panel for the security dashboard)
- [ ] Structured logging with scan lifecycle events

### v0.3 -- Smarts

- [ ] Image exclude patterns (skip known-noisy base images)
- [ ] Delta detection -- log/alert only on newly appeared CVEs between rescans
- [ ] Configurable severity threshold for alerts
- [ ] Retry backoff tuning

### Future Ideas

- [ ] SBOM generation and storage
- [ ] Integration with Renovate -- correlate "fix available" with pending Renovate PRs
- [ ] Webhook/notification on new critical CVEs (Discord, ntfy)
- [ ] Web UI for browsing scan results (or keep it Grafana-only?)
