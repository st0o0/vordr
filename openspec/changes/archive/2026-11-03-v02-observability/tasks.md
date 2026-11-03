## 1. Structured Logging

- [x] 1.1 Add `event` attribute to all log calls in `internal/watcher/watcher.go` (image_discovered, image_removed)
- [x] 1.2 Add `event` attribute to all log calls in `internal/worker/pool.go` (scan_started, scan_completed, scan_failed, scan_skipped)
- [x] 1.3 Add `event` attribute to all log calls in `internal/scheduler/scheduler.go` (scheduled_rescan)
- [x] 1.4 Add `event` attribute to startup/shutdown log in `cmd/vordr/main.go`

## 2. Vulnerability Types

- [x] 2.1 Add `Vulnerability` struct to `internal/scanner/types.go` (ID, PkgName, InstalledVersion, FixedVersion, Severity)
- [x] 2.2 Add `Vulnerabilities []Vulnerability` field to `ScanResult` in `internal/trivy/client.go`
- [x] 2.3 Update Trivy JSON parsing to populate Vulnerabilities slice alongside Counts/Fixable
- [x] 2.4 Update Trivy client tests with vulnerability detail assertions

## 3. Store Vulnerability Details

- [x] 3.1 Add `Vulnerabilities []scanner.Vulnerability` field to `ImageState` in `internal/scanner/types.go`
- [x] 3.2 Update `store.UpdateResult()` to accept and store vulnerabilities
- [x] 3.3 Update `store.Snapshot()` to copy vulnerabilities slice
- [x] 3.4 Update store tests for vulnerability details

## 4. Worker Integration

- [x] 4.1 Update `worker.Pool.scan()` to pass `result.Vulnerabilities` to `store.UpdateResult()`

## 5. Trivy DB Age Monitoring

- [x] 5.1 Create `internal/trivy/dbinfo.go` with `DBInfo` struct and `CheckDBVersion(ctx, cacheDir)` function that runs `trivy version --format json`
- [x] 5.2 Create `internal/trivy/dbcheck.go` with `DBChecker` struct — periodic goroutine (30min interval, runs once on start)
- [x] 5.3 Add `vordr_trivy_db_updated_timestamp_seconds` gauge to metrics collector (read from atomic pointer)
- [x] 5.4 Wire DBChecker into errgroup lifecycle in `cmd/vordr/main.go`
- [x] 5.5 Write tests for DB version JSON parsing

## 6. Details Endpoint

- [x] 6.1 Create `internal/metrics/details.go` with `/details` HTTP handler — reads Store.Snapshot(), returns JSON
- [x] 6.2 Support `?image=name:tag` query parameter for filtering
- [x] 6.3 Register `/details` route on the existing metrics server mux
- [x] 6.4 Write tests for details handler (JSON output, filtering, empty store)

## 7. Grafana Dashboard

- [x] 7.1 Create `grafana/vordr-dashboard.json` with panels: stat row (discovered/scanned/pending/errors), severity bar chart, per-image table, DB age gauge, scan duration history

## 8. Verification

- [x] 8.1 Run `go build ./...` and `go vet ./...`
- [x] 8.2 Run `go test ./...` — all tests pass
- [x] 8.3 Run `docker compose build` and `docker compose up` — verify /details and /metrics work
