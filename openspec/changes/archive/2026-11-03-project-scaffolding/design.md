## Context

Vordr ist ein neues Go-Projekt, das als CVE-Scanner für Homelab Docker Stacks dient. Es existiert bisher nur eine ROADMAP.md. Alle anderen Go-Repos des Autors (bifrost, eir, forseti, mjolnir, ran, var) folgen einem konsistenten Template mit identischer CI/CD-Infrastruktur, Tooling-Konfiguration und Projektstruktur. Vordr soll dieses Template übernehmen, um vom ersten Commit an vollständig funktionsfähig zu sein.

Das Scaffolding ist bewusst als eigener Change getrennt von der v0.1-Implementierung, damit die Infrastruktur unabhängig reviewbar ist und die Implementierung auf einer stabilen Basis aufbaut.

## Goals / Non-Goals

**Goals:**
- Vordr-Repo auf den gleichen Stand bringen wie die anderen Go-Repos
- CI/CD, Linting, Release-Automation, Dependency-Updates ab Tag 1
- cmd/vordr/ mit lauffähigem Binary (version + healthcheck Subcommands)
- internal/ Package-Struktur vorbereitet für v0.1-Implementierung
- Metrics-Server mit /metrics und /healthz auf :9550
- Config-Loading via Env-Vars nach bewährtem Pattern
- CLAUDE.md + AGENTS.md für AI-gestützte Entwicklung

**Non-Goals:**
- Keine Implementierung der Scanner-Logik (Watcher, Scheduler, WorkerPool, Store, Trivy-Client)
- Keine funktionalen Tests — die kommen mit der Implementierung
- Kein Grafana Dashboard oder Alloy-Config
- Kein compose.override für homelab-infra (kommt mit Deployment)

## Decisions

### D1: Projektstruktur nach bestehendem Template

**Entscheidung:** Identische Struktur wie bifrost/eir/forseti.

**Begründung:** Konsistenz über alle Repos. Jeder, der ein Repo kennt, findet sich sofort zurecht. CI-Workflows sind getestet und funktionieren mit dieser Struktur.

**Alternativen:**
- Flat-Layout ohne cmd/internal — abgelehnt, skaliert nicht wenn Packages wachsen
- Monorepo mit anderen Tools — abgelehnt, jedes Tool hat eigenen Release-Zyklus

### D2: Scratch als Docker Base-Image

**Entscheidung:** `scratch` (nicht alpine).

**Begründung:** Vordr hat keine externen Abhängigkeiten (kein Trivy-Binary, keine Shell nötig). CGO_ENABLED=0 produziert ein statisches Binary. Scratch minimiert Angriffsfläche und Image-Größe (~10MB).

**Alternativen:**
- Alpine — nur nötig wenn externe Binaries gebraucht werden (wie mjolnir für NUT)
- Distroless — unnötig, scratch reicht für statische Go-Binaries

### D3: Metrics Port :9550

**Entscheidung:** Standard-Port :9550 wie alle anderen Repos.

**Begründung:** Konsistenz. Alloy-Scrape-Config kann den gleichen Port für alle eigenen Services verwenden. Die ROADMAP erwähnte :8080, aber das ist der allgemeine HTTP-Port — :9550 ist der etablierte Metrics-Port im Homelab.

### D4: Stub-Packages in internal/

**Entscheidung:** Packages für watcher/, scheduler/, worker/, store/, trivy/ werden als leere Go-Files angelegt (nur `package <name>`).

**Begründung:** Definiert die Package-Struktur und Import-Pfade vorab. Die Implementierung füllt sie. config/ und metrics/ bekommen funktionalen Code, da sie für das lauffähige Binary gebraucht werden.

**Alternativen:**
- Keine Stubs — abgelehnt, dann fehlt die sichtbare Architektur im Repo
- Volle Interfaces definieren — abgelehnt, gehört zur Implementierung nicht zum Scaffolding

### D5: CLAUDE.md + AGENTS.md nach var/forseti Muster

**Entscheidung:** CLAUDE.md mit Projekt-Kontext, Architektur-Überblick, Build-Commands, Konventionen. AGENTS.md mit Code-Change-Instruktionen und OpenSpec-Nutzung.

**Begründung:** Ermöglicht AI-gestützte Entwicklung von Anfang an. Die Dateien dokumentieren das „Wie arbeiten wir hier" und nicht das „Was macht der Code" (das steht im Code selbst).

### D6: Dependencies

**Entscheidung:**
- `github.com/docker/docker` — Docker API Client
- `github.com/prometheus/client_golang` — Prometheus Metrics
- `golang.org/x/sync` — errgroup für strukturierte Concurrency

**Begründung:** Minimale Dependencies die für das MVP gebraucht werden. Kein HTTP-Framework — stdlib `net/http` reicht (wie in allen anderen Repos). Trivy-Client nutzt ebenfalls stdlib HTTP.

## Risks / Trade-offs

- **[Go-Version Drift]** → Repos nutzen teils unterschiedliche Minor-Versionen (1.26.0 bis 1.26.10). Vordr startet mit der neuesten stabilen Version. Renovate hält sie aktuell.
- **[Stub-Packages könnten veralten]** → Wenn die Implementierung die Architektur ändert, müssen Stubs angepasst werden. Risiko ist gering da die ROADMAP die Struktur klar vorgibt.
- **[github-workflows Abhängigkeit]** → Alle CI delegiert an st0o0/github-workflows. Wenn dieses Repo Probleme hat, sind alle Repos betroffen. Akzeptiertes Risiko — Zentralisierung ist hier gewollt.
