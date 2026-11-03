## Why

Das Repo hat jetzt die komplette Infrastruktur (CI, Tooling, Dockerfile, Metrics-Server), aber noch keine Scanner-Funktionalität. Dieser Change implementiert die v0.1 MVP-Logik: Container-Discovery, Trivy-Scanning, und Prometheus-Metriken — alles was nötig ist, damit vordr tatsächlich CVEs findet und reportet.

## What Changes

- **Store** (`internal/store/`): Thread-safe Image-State-Management mit sync.RWMutex, Snapshot für Metrics-Scrape, Add/Update/Remove Operationen
- **Watcher** (`internal/watcher/`): Docker API Container-Listing beim Start + Event-Stream für Start/Stop Events, Schreibt ImageRef auf scanCh
- **Scheduler** (`internal/scheduler/`): Periodic Rescan-Ticker (konfigurierbar, default 6h), schreibt bekannte Images auf scanCh
- **Worker Pool** (`internal/worker/`): N Goroutines lesen von scanCh, scannen via Trivy Client, updaten Store
- **Trivy Client** (`internal/trivy/`): HTTP Client für Trivy Server Scan-API, JSON Response Parsing, Severity-Counts Extraktion
- **Prometheus Metrics** (`internal/metrics/`): Erweiterung des bestehenden Stubs um alle vordr_* Metriken aus der ROADMAP (vulnerability counts, scan status, operational gauges)
- **Main Integration** (`cmd/vordr/main.go`): Erweiterung um errgroup-basierte Lifecycle aller Komponenten, scanCh Wiring, staggered Startup
- **Alloy Scrape Config**: Feature-Config für den bestehenden Alloy Stack

## Capabilities

### New Capabilities
- `container-discovery`: Docker API container listing + event stream, image reference extraction, container-to-image mapping
- `scan-orchestration`: Worker pool, scan channel, staggered startup scheduling, periodic rescan
- `trivy-client`: HTTP client for Trivy server scan API, JSON response parsing, vulnerability count extraction
- `image-state-store`: Thread-safe image state management, state machine (Pending→Scanning→Scanned/Error), snapshot for metrics
- `vulnerability-metrics`: All vordr_* Prometheus metrics — vulnerability counts, scan status, operational gauges, container-image info mapping

### Modified Capabilities
- `metrics-server`: Erweiterung des bestehenden Metrics-Stubs um tatsächliche vordr_* Collector-Registrierungen
- `config-loading`: Keine Spec-Änderung — Config-Struct und Env-Vars sind bereits definiert

## Impact

- **Neue Abhängigkeit**: `github.com/docker/docker` Client wird jetzt tatsächlich benutzt (war bisher nur in go.mod)
- **Docker Socket**: vordr braucht read-only Zugang zu `/var/run/docker.sock`
- **Trivy Server**: Muss erreichbar sein unter VORDR_TRIVY_ADDR
- **cmd/vordr/main.go**: Signifikante Erweiterung — von "start and wait" zu vollem errgroup-Lifecycle
- **internal/metrics/metrics.go**: Von Stub zu vollständiger Metrik-Registrierung
- **Alloy Stack**: Neues Feature-Config-File für den Scrape-Job
