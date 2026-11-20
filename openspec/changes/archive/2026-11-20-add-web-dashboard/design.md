## Context

Vordr is a single-binary CVE scanner for homelab Docker stacks. It already exposes an HTTP server (`internal/metrics/server.go`) with `/metrics`, `/healthz`, and `/details` endpoints. The `/details` endpoint returns a JSON array of all tracked images with their vulnerability details. There is no visual interface — users rely on Grafana dashboards or direct API calls.

## Goals / Non-Goals

**Goals:**
- Provide an at-a-glance vulnerability overview accessible via browser
- Zero additional dependencies — stdlib `embed`, `io/fs`, `net/http` only
- Single self-contained HTML file (no build pipeline, no framework)
- Dashboard consumes existing endpoints (same-origin, no CORS)

**Non-Goals:**
- No server-side rendering or templating
- No client-side routing or SPA framework
- No write operations (triggering scans, dismissing CVEs)
- No authentication (Vordr relies on network-level access control)
- No WebSocket or SSE for real-time updates

## Decisions

### 1. `embed.FS` with single HTML file

The dashboard is a single `dashboard.html` compiled into the binary via `//go:embed`. A new package `internal/ui/` owns the embed directive and exposes a `Handler() http.Handler` function.

**Why over `html/template`**: Client-side JS enables interactive filtering and expand/collapse without round-trips. A template approach would require JS anyway for interactivity, negating its advantage.

**Why single file over multi-file**: CSS and JS inline in one HTML file. No asset pipeline, no cache-busting, no relative path issues. The file is expected to be <50KB — well within reason for embedding.

### 2. Route on `GET /`

The root route is currently unregistered (returns 404). Serving the dashboard there is natural — Vordr is a scanner with a metrics endpoint, not a general API server. Users opening the host:port in a browser get the dashboard.

**Alternative considered**: `/ui` with redirect from `/`. Adds complexity for no real benefit since there are no competing uses for `/`.

### 3. `internal/ui/` package exposes `Handler()`

```
internal/ui/
  embed.go        — //go:embed dashboard.html, Handler() function
  dashboard.html  — self-contained HTML+CSS+JS
  handler_test.go — httptest verification
```

`metrics/server.go` imports `ui` and wires: `mux.Handle("GET /", ui.Handler())`. The `ui` package has no dependency on `metrics` — clean one-way import.

**Why not embed in metrics package directly**: Separation of concerns. The UI is a presentation layer, metrics is data. The `ui` package can grow independently (e.g., adding favicon.ico later) without touching metrics.

### 4. Auto-refresh via polling

`setInterval` fetching `/details` and `/healthz` every 30 seconds. Simple, predictable, no connection management. Scan intervals are measured in hours — 30s polling is more than sufficient.

### 5. Dashboard layout

- **Header**: vörðr logo/title, version, uptime, health status dot
- **Summary tiles**: Total vulns, Critical, High, Medium, Low counts
- **Image table**: Rows with image name, status, vuln count, last scan time
- **Expand/collapse**: Click image row to show vulnerability detail table
- **Severity badges**: Color-coded (Critical=red, High=orange, Medium=yellow, Low=blue)

## Risks / Trade-offs

- **[Binary size increase]** → Expected <50KB for the HTML file. Negligible compared to the ~200MB Trivy binary bundled in the container.
- **[No dark mode]** → Keep initial scope minimal. Can add via `prefers-color-scheme` media query later.
- **[Polling load]** → 2 fetches every 30s is trivial. The `/details` handler iterates the in-memory store with a read lock — microsecond response times.
- **[Go mux routing precedence]** → `GET /` in Go 1.22+ ServeMux is a catch-all for unmatched routes. Existing explicit routes (`/metrics`, `/healthz`, `/details`) take precedence — no conflict.
