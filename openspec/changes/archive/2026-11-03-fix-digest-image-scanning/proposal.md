## Why

Images die per Digest gepullt wurden (z.B. `ghcr.io/home-operations/radarr:rolling@sha256:52d09...`) können nicht gescannt werden. Trivy CLI gibt `FATAL: failed to analyze layer` zurück, weil die Digest-Referenz nicht im lokalen Docker-Store auflösbar ist. Im Live-Test scheiterten 3 von 4 Images — alle die per `@sha256` deployed waren.

Ursache: `ImageRef.String()` gibt die volle Referenz inkl. Digest zurück, die als Trivy-Scan-Target nicht funktioniert. Trivy braucht entweder einen Tag-basierten Namen (`nginx:latest`) oder eine Image-ID (`sha256:52d09...`).

## What Changes

- **ImageRef um ScanTarget erweitern**: Neues Feld `ImageID` das die Docker Image-ID (`sha256:...`) enthält — wird als Trivy-Scan-Target verwendet
- **ParseImageRef Digest-Handling**: `@sha256:...` korrekt erkennen und abspalten. Image-Name und Tag für Labels behalten, Digest separat
- **Watcher**: Container's `ImageID` aus Docker API extrahieren und in ImageRef speichern
- **Worker Pool**: `ref.ImageID` statt `ref.String()` als Trivy-Scan-Target nutzen (Fallback auf `ref.String()` wenn kein ImageID)
- **Store**: Key bleibt `ref.String()` (für Dedup), Scan nutzt `ref.ImageID`

## Capabilities

### New Capabilities

_Keine — bestehende Capability wird gefixt._

### Modified Capabilities

- `container-discovery`: ImageRef wird um ImageID erweitert, ParseImageRef handelt Digest-Referenzen
- `image-state-store`: Store-Key-Logik bleibt gleich, aber ImageRef hat neues Feld
- `scan-orchestration`: Worker nutzt ImageID als Scan-Target statt String()-Repräsentation

## Impact

- **internal/scanner/types.go**: ImageRef bekommt `ImageID` Feld
- **internal/watcher/watcher.go**: Extraktion von ImageID aus Docker Container-Metadata
- **internal/watcher/watcher_test.go**: Tests für Digest-Parsing
- **internal/worker/pool.go**: Scan-Target-Logik
- Keine Config-Änderung, keine Dockerfile-Änderung
