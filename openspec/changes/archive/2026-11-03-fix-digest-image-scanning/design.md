## Context

Beim Live-Test scheiterten 3 von 4 Images. Trivy CLI kann Images über zwei Wege scannen:
1. **Image-Name** (`nginx:latest`) — funktioniert nur wenn Docker den Namen kennt (nicht bei digest-only pulls)
2. **Image-ID** (`sha256:52d09...`) — funktioniert immer, greift direkt auf den lokalen Docker-Store zu

Docker Container API liefert beides:
- `container.Image` → die volle Referenz inkl. Digest (`ghcr.io/...@sha256:...`) — unbrauchbar für Trivy
- `container.ImageID` → die lokale Image-ID (`sha256:...`) — funktioniert immer

## Goals / Non-Goals

**Goals:**
- Alle lokal vorhandenen Images scannen, unabhängig davon wie sie gepullt wurden
- Prometheus-Labels behalten den menschenlesbaren Image-Namen (nicht die sha256-ID)
- Kein Breaking Change an der Store-Logik oder Metrics-Output

**Non-Goals:**
- Remote-Registry-Scanning (nur lokal vorhandene Images)
- OCI-Reference-Parsing Library einbinden

## Decisions

### D1: ImageRef bekommt ImageID Feld

```go
type ImageRef struct {
    Image   string // "grafana/grafana" — für Labels
    Tag     string // "13.3.0" oder "rolling" — für Labels
    ImageID string // "sha256:52d09..." — für Trivy Scan-Target
}
```

**Begründung:** Trennung von Display (Image+Tag für Prometheus) und Scan-Target (ImageID für Trivy). `String()` bleibt `Image:Tag` für Store-Key und Log-Output.

### D2: ScanTarget() Methode

```go
func (r ImageRef) ScanTarget() string {
    if r.ImageID != "" {
        return r.ImageID
    }
    return r.String()
}
```

**Begründung:** Fallback auf `Image:Tag` wenn kein ImageID verfügbar (z.B. in Tests). ImageID ist der bevorzugte Scan-Pfad.

### D3: ParseImageRef handelt drei Fälle

1. `nginx:1.25` → Image=`nginx`, Tag=`1.25`
2. `ghcr.io/org/app:rolling@sha256:abc...` → Image=`ghcr.io/org/app`, Tag=`rolling`
3. `nginx` → Image=`nginx`, Tag=`latest`

Der `@sha256:...` Teil wird bei Parsing abgeschnitten — er ist für Labels irrelevant. Die eigentliche Scan-ID kommt aus `container.ImageID`.

### D4: Watcher extrahiert ImageID

In der Container-List und Event-Handler:
- Container-List: `container.ImageID` direkt verfügbar
- Events: `msg.Actor.Attributes["imageID"]` (nicht verfügbar in events) → Fallback: `docker inspect` oder ImageID im Store nachschlagen

**Pragmatische Lösung:** Bei Events ist die ImageID nicht im Event-Payload. Aber neue Images werden nur bei `start` Events relevant, und dort können wir einen kurzen `docker inspect` machen. Oder besser: wir speichern die ImageID im Store bei Add() und nutzen sie von dort.

## Risks / Trade-offs

- **[docker inspect bei Events]** → Minimal overhead, passiert nur bei neuen Images. Alternativ: ImageID aus dem Store nachschlagen wenn das Image bereits bekannt ist.
- **[Tag-Parsing Heuristik]** → Einfaches String-Split auf `@` vor dem Tag-Parsing. Deckt die gängigen Fälle ab.
