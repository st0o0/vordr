## Context

v0.2 ist komplett: structured logging, /details, DB-age metric, Grafana Dashboard. Die vier v0.3 Features sind unabhängig voneinander und können parallel implementiert werden.

## Goals / Non-Goals

**Goals:**
- Images per Glob-Pattern von Scans ausschließen
- Bei Rescans nur neue CVEs loggen (Delta)
- Severity-Schwelle für "alert-worthy" Delta-CVEs
- Automatische Retry mit Backoff für fehlgeschlagene Scans

**Non-Goals:**
- Push-Notifications (Discord/ntfy) — Future
- Webhook-Integration — Future
- Per-Image Custom-Config — zu komplex für v0.3

## Decisions

### D1: Exclude Patterns — Glob via path.Match

`VORDR_EXCLUDE_IMAGES` ist eine Komma-separierte Liste von Glob-Patterns:
```
VORDR_EXCLUDE_IMAGES=aquasec/trivy:*,*devcontainer*,vordr-*
```

Matching via `path.Match` gegen `image:tag`. Check passiert im Watcher (bei Discovery) und im Scheduler (bei Rescan). Der Store weiß nichts von Excludes.

### D2: Delta Detection — Set-Differenz der CVE-IDs

Der Store bekommt ein `PreviousVulnIDs map[string]bool` pro Image. Nach jedem Scan:
1. Neue CVE-IDs = aktuelle IDs minus vorherige IDs
2. Falls neue CVEs vorhanden: `event=new_cves_detected` mit Count und Severity-Breakdown
3. Vorherige IDs werden mit aktuellen überschrieben

Delta-Log passiert im Worker nach dem Store-Update. Die Metrics zeigen weiterhin den vollen Zustand.

### D3: Severity Threshold — Ordinal Comparison

```go
var severityOrder = map[string]int{
    "UNKNOWN": 0, "LOW": 1, "MEDIUM": 2, "HIGH": 3, "CRITICAL": 4,
}
```

`VORDR_ALERT_SEVERITY` (default `HIGH`) bestimmt die Mindest-Severity für Delta-Alerts. Ein neues CVE mit Severity >= Threshold wird als `alert_worthy=true` geloggt.

### D4: Retry Backoff — Separate Retry-Queue

Ein neuer `internal/retry/queue.go`:
- Fehlgeschlagene Images werden in eine Retry-Queue geschoben (nach dem Scan-Error im Worker)
- Queue-Goroutine wartet den Backoff ab, dann schiebt das Image zurück auf scanCh
- Backoff: 30s → 1m → 2m → 4m → 8m → 15m (max), reset nach erfolgreichem Scan
- `VORDR_RETRY_MAX_BACKOFF` (default `15m`)

```
Worker ──scan fails──▶ RetryQueue ──backoff──▶ scanCh ──▶ Worker
```

Der Scheduler queued weiterhin alle Images alle 6h. Retry ist zusätzlich und unabhängig.

### D5: Config-Erweiterung

```go
type Config struct {
    // ... existing ...
    ExcludeImages   []string      // glob patterns
    AlertSeverity   string        // "HIGH" default
    RetryMaxBackoff time.Duration // 15m default
}
```

Neue Env-Vars:
- `VORDR_EXCLUDE_IMAGES` (default leer)
- `VORDR_ALERT_SEVERITY` (default `HIGH`)
- `VORDR_RETRY_MAX_BACKOFF` (default `15m`)

## Risks / Trade-offs

- **[Retry-Queue Memory]** → Bounded: max 1 Entry pro fehlgeschlagenes Image. Bei 50 Images max 50 Entries.
- **[Delta Detection false positives]** → Beim ersten Scan sind alle CVEs "neu". Erster Scan emittiert kein Delta-Event (kein Previous).
- **[Glob Patterns Limitierung]** → `path.Match` unterstützt `*`, `?`, `[...]`. Keine Regex, keine `**`. Reicht für Homelab-Scope.
