## Why

v0.1 scannt Images und exponiert Severity-Counts als Prometheus-Metriken. Aber: Debugging ist schwierig (welche CVEs genau?), die Trivy DB Freshness ist unsichtbar, es gibt kein Dashboard, und das Logging ist inkonsistent. v0.2 macht vordr production-ready für ein Homelab-Monitoring-Setup.

## What Changes

- **Structured Logging**: Konsistente slog-Events mit Event-Labels für den gesamten Scan-Lifecycle (image_discovered, scan_started, scan_completed, scan_failed, image_removed, scheduled_rescan). Alle Events als Key-Value-Pairs.
- **Trivy DB Age Metric**: `vordr_trivy_db_updated_timestamp_seconds` Gauge — periodische Abfrage via `trivy version --format json`, Alert-fähig wenn DB älter als 48h.
- **`/details` JSON Endpoint**: GET /details — gibt pro Image die vollständige CVE-Liste zurück (VulnerabilityID, Package, InstalledVersion, FixedVersion, Severity). Trivy-Client muss volle CVE-Daten speichern statt nur Counts.
- **Grafana Dashboard**: JSON-Dashboard-Definition für die vordr_* Metriken — CVE-Übersicht, Severity-Verteilung, Scan-Status, DB-Alter.

## Capabilities

### New Capabilities
- `details-endpoint`: HTTP GET /details JSON API mit per-image CVE-Listen
- `trivy-db-monitoring`: Trivy DB age metric via `trivy version` CLI
- `grafana-dashboard`: JSON dashboard definition for vordr metrics

### Modified Capabilities
- `trivy-client`: ScanResult erweitert um vollständige Vulnerability-Details (nicht nur Counts)
- `image-state-store`: ImageState speichert Vulnerability-Details für /details Endpoint
- `metrics-server`: Neuer /details Route + Trivy DB age metric im Collector
- `scan-orchestration`: Structured logging Events in allen Komponenten

## Impact

- **internal/trivy/**: ScanResult bekommt Vulnerabilities-Liste
- **internal/scanner/types.go**: Neuer Vulnerability-Typ
- **internal/store/**: ImageState speichert Vulnerability-Details
- **internal/metrics/**: DB-Age Gauge, /details Handler
- **internal/watcher/**, **internal/scheduler/**, **internal/worker/**: Structured log events
- **Neues File**: grafana/vordr-dashboard.json
- **Keine Config-Änderung**, keine Dockerfile-Änderung
