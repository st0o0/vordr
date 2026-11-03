## Context

Vordr's `internal/trivy/client.go` sendet HTTP POST an Trivy Server's Twirp-Endpoint. Dieser Endpoint ist ein Cache-Backend, kein Scan-Service — er erwartet vorab hochgeladene Layer-Blobs. Die Research hat gezeigt: das offizielle Trivy Operator Projekt nutzt die Trivy CLI als Subprocess. Vordr soll das gleiche tun.

Bestehender Code der sich ändert:
- `internal/trivy/client.go` — komplett neu
- `internal/trivy/client_test.go` — komplett neu
- `internal/config/config.go` — TrivyAddr → TrivyCacheDir + TrivyTimeout
- `cmd/vordr/main.go` — Client-Konstruktor
- `Dockerfile` — neues Base-Image
- `docker-compose.yml` — ein Service statt zwei

## Goals / Non-Goals

**Goals:**
- Funktionierende Trivy-Scans via CLI Subprocess
- Ein-Container-Architektur
- Trivy DB Cache auf Volume für Persistenz
- Gleiche ScanResult-Schnittstelle für Worker Pool (keine Änderung an worker/pool.go)

**Non-Goals:**
- Trivy Server-Mode Support (kann später als Option zurückkommen)
- Trivy als Go-Library einbinden
- Trivy DB Update Monitoring (v0.2)

## Decisions

### D1: exec.CommandContext mit JSON-Output

**Entscheidung:**
```go
cmd := exec.CommandContext(ctx, "trivy", "image",
    "--format", "json",
    "--scanners", "vuln",
    "--cache-dir", cacheDir,
    "--quiet",
    imageRef)
```

**Begründung:** `--format json` gibt strukturiertes Output auf stdout. `--scanners vuln` beschränkt auf Vulnerability-Scanning (kein Secret/Config-Scan). `--quiet` unterdrückt Progress-Output auf stderr. `--cache-dir` erlaubt persistenten DB-Cache.

**JSON-Parsing:** Minimale Go-Structs die nur `Results[].Vulnerabilities[].Severity` und `FixedVersion` lesen. Keine Dependency auf Trivy's Go-Types.

### D2: Alpine statt scratch als Base-Image

**Entscheidung:** Final-Stage basiert auf Alpine (nicht scratch).

**Begründung:** Trivy braucht:
- Eine Shell (für einige interne Operations)
- ca-certificates (für HTTPS DB-Downloads)
- Grundlegende C-Library (Trivy ist ein Go-Binary, aber nutzt CGO für einige Features)

Alpine ist der kleinste Kompromiss. Image-Größe: ~200MB (Trivy ~180MB + vordr ~10MB + Alpine ~5MB).

**Alternative:** Distroless — abgelehnt, Trivy braucht Shell-Zugriff für einige DB-Operationen.

### D3: Multi-Stage Dockerfile mit Trivy-Binary COPY

**Entscheidung:**
```dockerfile
FROM aquasec/trivy:0.76.1 AS trivy
FROM golang:1.27-alpine AS build
# ... build vordr ...
FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=trivy /usr/local/bin/trivy /usr/local/bin/trivy
COPY --from=build /vordr /vordr
```

**Begründung:** Pinned Trivy-Version (nicht `latest`), Binary wird aus dem offiziellen Image extrahiert. Kein separater Download nötig. Renovate kann die Trivy-Version automatisch updaten.

### D4: Config-Änderung — VORDR_TRIVY_ADDR → VORDR_TRIVY_CACHE_DIR

**Entscheidung:**
- `VORDR_TRIVY_ADDR` entfällt (kein Server mehr)
- `VORDR_TRIVY_CACHE_DIR` (default `/var/cache/trivy`) — wo Trivy die DB cached
- `VORDR_TRIVY_TIMEOUT` (default `5m`) — Timeout pro Scan

**Begründung:** Die DB muss zwischen Restarts persistent sein. Docker Volume auf `VORDR_TRIVY_CACHE_DIR` verhindert Re-Downloads. Timeout schützt vor hängenden Scans.

### D5: Stderr-Handling

**Entscheidung:** Trivy-Stderr wird als Debug-Log geloggt (nicht als Error). Exit-Code 0 = Erfolg, alles andere = Error.

**Begründung:** Trivy schreibt Progress und Warnings auf stderr, auch bei erfolgreichen Scans. Das ist kein Fehler.

### D6: DB-Update Strategie

**Entscheidung:** Trivy updated die DB automatisch beim ersten Scan (Standard-Verhalten). Kein `--skip-db-update` Flag.

**Begründung:** Für Homelab ist der automatische Download beim Start ausreichend. DB wird auf dem Cache-Volume persistent — nach dem ersten Download sind Restarts schnell.

## Risks / Trade-offs

- **[Image-Größe steigt von ~10MB auf ~200MB]** → Akzeptiert. Funktionalität > Größe. Trivy-Binary ist der dominierende Faktor.
- **[Erster Scan dauert länger (DB-Download)]** → ~30s für die initiale DB. Danach gecached. Volume-Persistenz vermeidet Re-Downloads.
- **[Trivy CLI Stabilität]** → JSON-Output Format ist seit SchemaVersion 2 stabil. Minimale Structs reduzieren Breaking-Change-Risiko.
- **[Kein Trivy Server = kein shared Cache]** → Für ein Homelab mit einem vordr-Instance kein Problem. Falls mehrere Instances nötig werden, kann Server-Mode als Option zurückkommen.
