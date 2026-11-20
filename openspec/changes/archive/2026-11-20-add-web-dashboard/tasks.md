## 1. UI Package Setup

- [x] 1.1 Create `internal/ui/embed.go` with `//go:embed dashboard.html` and `Handler()` function returning `http.Handler`
- [x] 1.2 Create `internal/ui/handler_test.go` — verify Handler returns 200, Content-Type text/html, and non-empty body

## 2. Dashboard HTML

- [x] 2.1 Create `internal/ui/dashboard.html` — self-contained HTML+CSS+JS file with header (title, version, uptime, health dot), severity summary tiles, and image list
- [x] 2.2 Implement `/details` fetch and render: image rows with name, tag, status, vuln count, relative scan time
- [x] 2.3 Implement expand/collapse per image row showing vulnerability detail table (CVE ID, package, installed, fixed, severity with color badges)
- [x] 2.4 Implement `/healthz` fetch and render version + uptime in header
- [x] 2.5 Implement auto-refresh every 30s with last-refresh timestamp display
- [x] 2.6 Handle empty state (no images, no vulnerabilities)

## 3. Server Integration

- [x] 3.1 Wire `ui.Handler()` into `metrics/server.go` on `GET /` route
- [x] 3.2 Update `NewServer` signature if needed to accept ui handler (or import directly)

## 4. Verification

- [x] 4.1 Run `go build ./...` and `go vet ./...` — ensure embed compiles cleanly
- [x] 4.2 Run `go test ./...` — all existing + new tests pass
- [x] 4.3 Run `golangci-lint run` — no new lint issues
