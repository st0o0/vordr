## Why

Vordr exposes vulnerability data via `/details` JSON and `/metrics` Prometheus endpoints, but viewing scan results requires either Grafana or curl. An embedded web dashboard gives immediate visibility into scan status without external tooling — open the browser, see your images and CVEs. This fits the single-binary homelab philosophy: no extra containers, no frontend build pipeline.

## What Changes

- Add a self-contained HTML+CSS+JS dashboard compiled into the binary via `go:embed`
- Serve the dashboard on `GET /` (root route, currently unregistered)
- Dashboard consumes existing `/details` and `/healthz` endpoints (same-origin)
- Displays severity summary tiles, image list with expandable per-image vulnerability details
- Auto-refreshes every 30 seconds
- Shows version and uptime from `/healthz` in the header

## Capabilities

### New Capabilities
- `web-dashboard`: Embedded HTML dashboard served from the Vordr binary, showing vulnerability scan results with severity summaries, per-image drill-down, and auto-refresh

### Modified Capabilities
- `metrics-server`: The HTTP server gains a `GET /` route serving the embedded dashboard

## Impact

- **New package**: `internal/ui/` with embed directive, handler, and HTML file
- **Modified**: `internal/metrics/server.go` — new route registration, new dependency on `internal/ui`
- **Binary size**: Increases by the size of the HTML file (expected <50KB)
- **No new Go dependencies**: Uses only `embed`, `io/fs`, `net/http` from stdlib
- **No breaking changes**: All existing endpoints unchanged
