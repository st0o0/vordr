## Context

Vordr hat ein vollständig scaffoldetes Repo mit CI, Metrics-Server-Stub, Config-Loading, und Healthcheck. Die ROADMAP.md definiert die Architektur im Detail: Watcher + Scheduler als Producers, WorkerPool als Consumer, Store als zentrale State-Verwaltung, HTTP-Server für Metrics. Dieses Design beschreibt, wie die einzelnen Komponenten zusammengesteckt werden.

Bestehender Code der erweitert wird:
- `cmd/vordr/main.go` — aktuell: config → metrics server → `<-ctx.Done()`. Wird zu: config → alle Komponenten via errgroup
- `internal/metrics/metrics.go` — aktuell: leerer Stub mit Registry-Ref. Wird zu: alle vordr_* Prometheus-Metriken
- `internal/config/config.go` — bleibt unverändert, Config-Struct hat bereits alle nötigen Felder

## Goals / Non-Goals

**Goals:**
- Alle laufenden Docker-Container entdecken und deren Images via Trivy scannen
- Vulnerability-Counts als Prometheus-Metriken exponieren
- Reagieren auf Container Start/Stop Events in Echtzeit
- Periodische Rescans für neue CVE-Disclosures
- Gestaffelter Startup-Scan um Trivy nicht zu überlasten
- Sauberes Shutdown aller Goroutines via Context + errgroup

**Non-Goals:**
- /details JSON Endpoint (v0.2)
- Trivy DB Age Tracking (v0.2)
- Grafana Dashboard (v0.2)
- Image Exclude Patterns (v0.3)
- Delta Detection (v0.3)
- SBOM Generation (Future)

## Decisions

### D1: Single scanCh shared between Producers and Consumers

**Entscheidung:** Ein einzelner `chan ImageRef` verbindet Watcher + Scheduler (Producers) mit dem WorkerPool (Consumer).

```
Watcher ──┐
           ├──▶ scanCh (buffered) ──▶ WorkerPool (N goroutines)
Scheduler ─┘                              │
                                    Store.Update()
```

**Begründung:** Einfachstes mögliches Design. Der Channel fungiert als natürlicher Backpressure-Mechanismus. Wenn alle Worker beschäftigt sind, blockieren die Producers bis ein Worker frei wird.

**Channel-Kapazität:** `2 * scanWorkers` — genug Buffer für einen kleinen Burst ohne Memory-Risiko.

**Alternativen:**
- Separate Channels pro Producer mit Fan-in — unnötige Komplexität
- Priority Queue (Events vor Scheduled) — Overhead nicht gerechtfertigt für Homelab-Scale

### D2: ImageRef als Channel-Typ

**Entscheidung:**
```go
type ImageRef struct {
    Image string // e.g. "grafana/grafana:13.3.0"
    Tag   string // e.g. "13.3.0" (parsed)
}
```

**Begründung:** Workers brauchen nur den Image-Identifier. Container-Mapping lebt im Store. Duplicate-Sends sind harmlos — der Worker prüft im Store ob ein Scan nötig ist.

### D3: Store als sync.RWMutex + map

**Entscheidung:** `sync.RWMutex` mit `map[string]*ImageState`. Key ist die Image-Referenz (ohne Tag-Normalisierung).

```go
type ImageState struct {
    Status     Status
    Containers []string
    Counts     map[string]int  // severity → count
    Fixable    map[string]int  // severity → fixable count
    ScannedAt  time.Time
    Duration   time.Duration
    Error      string
}
```

**Begründung:** Access-Pattern ist simpel: seltene Writes (nach jedem Scan, alle paar Minuten), häufige Reads (jeder Alloy-Scrape, alle 30s). RWMutex ist das richtige Primitive. `Snapshot()` kopiert den State für den Metrics-Handler — kein Lock während der HTTP-Response.

### D4: Watcher-Design — List + Watch Pattern

**Entscheidung:**
1. **Startup:** `docker container list` → alle laufenden Container → Images extrahieren → Store.Add() + scanCh
2. **Laufend:** `docker events --filter type=container --filter event=start --filter event=die` → Store.Add/Remove + scanCh bei neuen Images

**Begründung:** List+Watch ist das Standard-Pattern für eventuelle Konsistenz. Die initiale Liste fängt den Ist-Zustand, Events halten ihn aktuell. Docker Events API ist ein Long-Poll der nie terminiert — passt perfekt als Goroutine mit Context-Cancellation.

**Container-zu-Image-Extraktion:** `container.Image` + `container.ImageID` aus der Docker API. Tag-Parsing via String-Split (kein OCI-Library nötig für Homelab-Scope).

### D5: Staggered Startup

**Entscheidung:** Nach dem initialen Container-Listing werden Scan-Requests nicht sofort gequeued, sondern mit `ScanStagger` Abstand (default 30s) zeitversetzt gesendet.

```go
for i, ref := range initialImages {
    if i > 0 {
        select {
        case <-ctx.Done():
            return
        case <-time.After(cfg.ScanStagger):
        }
    }
    scanCh <- ref
}
```

**Begründung:** Trivy Server cached intern, aber der erste Scan pro Image dauert 10-30s (Download + Analyse). 40 Images * 30s = 20 Minuten bis alle initial gescannt sind — akzeptabel für ein Homelab.

### D6: Trivy Client — Stdlib HTTP

**Entscheidung:** `net/http` Client mit dem Trivy Server Scan-Endpoint. Request: POST mit Image-Ref. Response: JSON mit Vulnerability-Liste.

```
POST /twirp/trivy.scanner.v1.Scanner/Scan
Content-Type: application/json
{"target": "grafana/grafana:13.3.0", "artifact_type": 2}

→ Response: JSON mit Results[].Vulnerabilities[]
```

**Begründung:** Trivy Server hat ein simples JSON API. Kein gRPC-Client nötig — die Twirp-Endpoints akzeptieren plain JSON via HTTP POST. Stdlib `net/http` mit Timeout reicht.

**Fallback bei Scan-Fehler:** State → Error, Error-Message im Store, Worker loggt und macht mit nächstem Image weiter. Kein Retry im Worker — der Scheduler queued nach dem nächsten Interval automatisch.

### D7: Prometheus Metrics — Collect-on-Scrape via Custom Collector

**Entscheidung:** Ein Custom `prometheus.Collector` der bei jedem Scrape `Store.Snapshot()` aufruft und dynamisch Metriken generiert.

**Begründung:** Images kommen und gehen. Mit festen Labels würden wir stale Metriken akkumulieren. Ein Custom Collector gibt bei jedem Scrape nur die aktuell bekannten Images zurück — Prometheus sieht sofort wenn ein Image verschwindet.

**Metriken (aus ROADMAP):**
- `vordr_image_scan_status` (Gauge, labels: image, tag)
- `vordr_vulnerability_count` (Gauge, labels: image, tag, severity)
- `vordr_fixable_count` (Gauge, labels: image, tag, severity)
- `vordr_scan_timestamp_seconds` (Gauge, labels: image, tag)
- `vordr_scan_duration_seconds` (Gauge, labels: image, tag)
- `vordr_container_image_info` (Gauge, labels: container, image, tag)
- `vordr_images_discovered` (Gauge)
- `vordr_images_pending` (Gauge)
- `vordr_images_scanned` (Gauge)
- `vordr_images_errors` (Gauge)
- `vordr_scan_workers_active` (Gauge)
- `vordr_scan_queue_length` (Gauge)

### D8: Errgroup Lifecycle in main.go

**Entscheidung:**
```go
g, ctx := errgroup.WithContext(ctx)
g.Go(func() error { return metricsServer.Run(ctx) })
g.Go(func() error { return watcher.Run(ctx, scanCh) })
g.Go(func() error { return scheduler.Run(ctx, scanCh) })
g.Go(func() error { return workerPool.Run(ctx) })
if err := g.Wait(); err != nil { ... }
```

**Begründung:** Identisch zum Pattern in der ROADMAP. Alle Komponenten teilen denselben Context. Wenn einer fehlschlägt, wird der Context gecancelled und alles fährt sauber runter.

**ScanCh Ownership:** main.go erstellt den Channel und gibt ihn an Producer (Watcher, Scheduler) und Consumer (WorkerPool). main.go schließt den Channel nicht — die Context-Cancellation reicht für sauberes Shutdown.

## Risks / Trade-offs

- **[Trivy Server muss laufen]** → Wenn Trivy nicht erreichbar ist, gehen alle Scans in Error-State. Kein Retry-Backoff im MVP — der Scheduler queued beim nächsten Interval erneut. Akzeptabel.
- **[Docker Socket Zugang]** → Read-only Mount ist Pflicht. Vordr liest nur Container-Listings und Events, schreibt nie. Trotzdem: Docker Socket = Root-Äquivalent auf dem Host.
- **[Kein Dedup im scanCh]** → Watcher und Scheduler könnten das gleiche Image gleichzeitig queuen. Worker prüft Store-State und skippt wenn bereits Scanning. Minimal verschwendete CPU, kein Problem.
- **[Twirp API Stabilität]** → Trivy's interne API ist nicht offiziell stabil. Bei Trivy-Updates könnte sich das JSON-Format ändern. Risiko ist gering für Homelab und durch Image-Pinning kontrollierbar.
- **[Tag-Parsing Heuristik]** → Simple String-Split statt OCI-Reference-Parsing. Könnte bei exotischen Image-Referenzen (Registries mit Ports, Digest-only) fehlschlagen. Für Homelab-Docker-Compose ausreichend.
