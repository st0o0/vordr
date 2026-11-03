## Why

Der bisherige Trivy HTTP-Client funktioniert nicht — Trivy Server ist ein Cache-Backend, kein Scan-Service. Research (10 parallele Agents, Web + Source-Code-Analyse) hat gezeigt: alle ernstzunehmenden Projekte (Trivy Operator, Harbor Scanner, GitLab) nutzen die Trivy CLI als Subprocess. Vordr muss das gleiche tun. Gleichzeitig vereinfacht sich die Architektur: statt zwei Container (vordr + trivy server) reicht ein einzelner Container der das Trivy-Binary bundelt.

## What Changes

- **Trivy Client umschreiben** (`internal/trivy/`): HTTP-Client durch `exec.CommandContext("trivy", "image", "--format", "json", ...)` ersetzen. Gleiche `ScanResult`-Rückgabe, komplett andere Implementierung.
- **Dockerfile**: Von `scratch` auf multi-stage mit Trivy-Binary. Stage 1 baut vordr (wie bisher), Stage 2 kopiert Trivy-Binary von `aquasec/trivy`, finale Stage basiert auf Alpine (braucht Shell für Trivy + ca-certificates).
- **Config**: `VORDR_TRIVY_ADDR` entfernen, `VORDR_TRIVY_CACHE_DIR` hinzufügen (default `/var/cache/trivy`). Optional `VORDR_TRIVY_TIMEOUT` für CLI-Timeout (default 5m).
- **docker-compose.yml**: Von 2 Services (vordr + trivy) auf 1 Service (vordr) mit trivy-cache Volume.
- **Tests**: Trivy-Client-Tests von HTTP-Mock auf exec-Mock umstellen.
- **main.go**: `trivy.NewClient(cfg.TrivyAddr)` → `trivy.NewCLI(cfg.TrivyCacheDir, cfg.TrivyTimeout)`

## Capabilities

### New Capabilities

_Keine neuen Capabilities — die Scan-Fähigkeit existiert bereits, nur die Implementierung ändert sich._

### Modified Capabilities

- `trivy-client`: Von HTTP-Client zu CLI-Subprocess. Gleiche Scan-Semantik, andere Transport-Schicht.
- `config-loading`: `VORDR_TRIVY_ADDR` wird durch `VORDR_TRIVY_CACHE_DIR` und `VORDR_TRIVY_TIMEOUT` ersetzt.
- `deployment`: Von 2-Container-Setup (vordr + trivy server) auf 1-Container mit gebundeltem Trivy.
- `build-and-ci`: Dockerfile ändert sich von scratch auf Alpine mit Trivy-Binary.

## Impact

- **Dockerfile**: Signifikanter Umbau — neues Base-Image, zusätzliche Build-Stage
- **Image-Größe**: Von ~10MB (scratch) auf ~200MB (Alpine + Trivy + vordr)
- **internal/trivy/**: Kompletter Rewrite (client.go + client_test.go)
- **internal/config/**: Config-Struct ändert sich (TrivyAddr → TrivyCacheDir + TrivyTimeout)
- **cmd/vordr/main.go**: Client-Konstruktor-Aufruf ändert sich
- **docker-compose.yml**: Trivy-Service entfällt komplett
- **Keine Go-Dependency-Änderung**: Trivy wird nicht als Go-Module importiert, nur als Binary
