## Context

v0.1 ist komplett: Watcher, Worker Pool, Scheduler, Store, Trivy CLI, Prometheus Metrics. Logging existiert via slog, aber Events sind nicht konsistent strukturiert. Der /metrics Endpoint zeigt nur Counts, keine CVE-Details. Die Trivy DB Freshness ist nicht überwacht.

## Goals / Non-Goals

**Goals:**
- Konsistente, maschinenlesbare Log-Events für den gesamten Scan-Lifecycle
- Trivy DB Age als Prometheus-Metrik für Alerting
- Per-Image CVE-Details über HTTP abrufbar
- Grafana Dashboard als JSON-Definition (importierbar)

**Non-Goals:**
- Web UI für CVE-Browsing (Future)
- Push-Notifications bei neuen CVEs (v0.3+)
- Delta-Detection zwischen Scans (v0.3)

## Decisions

### D1: Structured Logging — Event-basierte slog Attributes

Alle Log-Events bekommen ein `event` Attribute:
```go
logger.Info("scan completed", "event", "scan_completed", "image", key, "duration", d, "counts", counts)
```

Events: `image_discovered`, `image_removed`, `scan_started`, `scan_completed`, `scan_failed`, `scan_skipped`, `scheduled_rescan`, `db_check`.

**Begründung:** Maschinenlesbar für Log-Aggregation (Loki/Alloy). `event` als Label erlaubt einfaches Filtern.

### D2: Trivy DB Age — Periodischer `trivy version` Check

Neuer Goroutine im errgroup-Lifecycle: alle 30 Minuten `trivy version --format json` ausführen und `vordr_trivy_db_updated_timestamp_seconds` aktualisieren.

```go
type DBInfo struct {
    UpdatedAt time.Time
}
```

Das Trivy-version-JSON enthält `VulnerabilityDB.UpdatedAt`. Atomic-Pointer für thread-safe Updates, Collector liest beim Scrape.

**Begründung:** Separater Goroutine statt pro-Scan Check — DB-Version ändert sich selten (1x/Tag), kein Overhead nötig.

### D3: Vulnerability Details — ScanResult erweitern

```go
type Vulnerability struct {
    ID              string
    PkgName         string
    InstalledVersion string
    FixedVersion    string
    Severity        string
}

type ScanResult struct {
    Counts          map[string]int
    Fixable         map[string]int
    Vulnerabilities []Vulnerability  // NEU
}
```

Der Trivy-Client parsed bereits `Results[].Vulnerabilities[]` — jetzt werden alle Felder behalten statt nur Severity/FixedVersion.

ImageState im Store bekommt ein `Vulnerabilities []Vulnerability` Feld. Snapshot() kopiert es mit.

### D4: /details Endpoint

```
GET /details
→ JSON Array von Image-Objekten mit CVE-Listen

GET /details?image=nginx:latest
→ Einzelnes Image-Objekt
```

Response-Format:
```json
[{
  "image": "nginx",
  "tag": "latest",
  "status": "scanned",
  "scanned_at": "2026-11-03T14:20:20Z",
  "vulnerabilities": [
    {"id": "CVE-2024-1234", "package": "libssl", "installed": "1.1.1", "fixed": "1.1.2", "severity": "HIGH"}
  ]
}]
```

Registriert auf dem bestehenden Metrics-Server mux (Port :9550).

### D5: Grafana Dashboard — JSON Provisioning

Ein statisches JSON-File `grafana/vordr-dashboard.json` das via Grafana Provisioning oder manuellen Import nutzbar ist. Panels:
- Stat: Total Images / Scanned / Pending / Errors
- Bar Chart: Vulnerability Counts by Severity (stacked)
- Table: Per-Image CVE Details
- Gauge: Trivy DB Age
- Timeline: Scan Activity

**Begründung:** JSON-File statt Dashboard-as-Code (Grafonnet) — simpler, direkt importierbar.

## Risks / Trade-offs

- **[Vulnerability-Details im Memory]** → Bei 50 Images mit je 100 CVEs = ~5000 Entries im Store. Für Homelab-Scale kein Problem. Bei größeren Setups könnte Memory relevant werden.
- **[/details kein Auth]** → Für Homelab-LAN akzeptabel. Keine sensiblen Daten (nur CVE-IDs und Versionen).
- **[Grafana Dashboard Maintenance]** → Manuell gepflegtes JSON. Bei Metric-Änderungen muss das Dashboard angepasst werden.
