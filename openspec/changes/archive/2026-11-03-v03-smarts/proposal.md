## Why

v0.2 liefert volle Observability, aber vordr ist noch nicht smart genug für den täglichen Betrieb: jedes Image wird gescannt (auch known-noisy base images), Rescans melden alle CVEs erneut statt nur Deltas, es gibt keinen Severity-Filter für Alerts, und fehlgeschlagene Scans werden erst beim nächsten 6h-Intervall wiederholt. v0.3 macht vordr intelligent genug dass man es einrichten und vergessen kann.

## What Changes

- **Image Exclude Patterns**: `VORDR_EXCLUDE_IMAGES` Env-Var mit Glob-Patterns (z.B. `aquasec/trivy:*,*devcontainer*`). Watcher und Scheduler überspringen gematchte Images.
- **Delta Detection**: Bei Rescans vergleicht der Worker die neuen CVE-IDs mit den vorherigen. Nur neu aufgetauchte CVEs werden geloggt. Neues `event=new_cves_detected` Log-Event.
- **Severity Threshold**: `VORDR_ALERT_SEVERITY` (default `HIGH`) — nur CVEs ab dieser Severity werden im Delta-Log als "alert-worthy" markiert. Beeinflusst nicht die Metrics (die zeigen immer alles).
- **Retry Backoff**: Fehlgeschlagene Scans werden nach exponential Backoff (30s, 1m, 2m, 4m, max 15m) automatisch erneut versucht statt auf das 6h-Rescan-Intervall zu warten.

## Capabilities

### New Capabilities
- `image-exclusion`: Glob-Pattern-basiertes Image-Filtering
- `delta-detection`: CVE-Delta-Erkennung zwischen Rescans mit Severity-basiertem Alerting
- `retry-backoff`: Exponential Backoff für fehlgeschlagene Scans

### Modified Capabilities
- `config-loading`: Neue Env-Vars VORDR_EXCLUDE_IMAGES, VORDR_ALERT_SEVERITY, VORDR_RETRY_MAX_BACKOFF
- `container-discovery`: Watcher filtert excluded Images
- `scan-orchestration`: Worker macht Delta-Detection, Scheduler respektiert Excludes
- `image-state-store`: Store tracked vorherige CVE-IDs für Delta-Detection

## Impact

- **internal/config/**: 3 neue Config-Felder
- **internal/watcher/**: Exclude-Check bei Discovery
- **internal/scheduler/**: Exclude-Check bei Rescan
- **internal/worker/**: Delta-Detection-Logik, Retry-Queue
- **internal/store/**: PreviousVulnIDs tracking
- **internal/scanner/types.go**: Keine Änderung nötig
- **cmd/vordr/main.go**: Retry-Queue Wiring
