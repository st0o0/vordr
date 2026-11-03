## 1. Shared Types

- [x] 1.1 Create `internal/scanner/types.go` with ImageRef (Image, Tag) and Status enum (Pending, Scanning, Scanned, Error)
- [x] 1.2 Create `internal/scanner/types.go` with ImageState struct (Status, Containers, Counts, Fixable, ScannedAt, Duration, Error)

## 2. Image State Store

- [x] 2.1 Implement `internal/store/store.go` with Store struct (sync.RWMutex + map[string]*ImageState)
- [x] 2.2 Implement Add(image ImageRef, container string) — adds image if new (Pending), appends container to list
- [x] 2.3 Implement Remove(container string) — removes container from image, deletes image if last container
- [x] 2.4 Implement SetScanning(image string) bool — transitions to Scanning, returns false if already Scanning
- [x] 2.5 Implement UpdateResult(image string, counts, fixable map[string]int, duration time.Duration) — sets Scanned + results
- [x] 2.6 Implement SetError(image string, errMsg string) — sets Error state
- [x] 2.7 Implement Snapshot() []ImageState — RLock, copy all state, RUnlock
- [x] 2.8 Implement Images() []ImageRef — returns all tracked image refs for scheduler
- [x] 2.9 Write unit tests for Store (concurrent add/remove, state transitions, snapshot isolation)

## 3. Trivy Client

- [x] 3.1 Implement `internal/trivy/client.go` with Client struct (baseURL, httpClient with 5min timeout)
- [x] 3.2 Implement Scan(ctx, imageRef string) (*ScanResult, error) — POST to Trivy server, parse JSON response
- [x] 3.3 Define ScanResult with Counts and Fixable maps (severity → count), extracted from Trivy's Results[].Vulnerabilities[]
- [x] 3.4 Write unit tests for Trivy client JSON parsing (mock HTTP responses)

## 4. Container Watcher

- [x] 4.1 Implement `internal/watcher/watcher.go` with Watcher struct (dockerClient, store, scanCh, stagger, logger)
- [x] 4.2 Implement Run(ctx) — initial container listing with staggered scan scheduling
- [x] 4.3 Implement event loop — subscribe to Docker container start/die events, update store, send new images to scanCh
- [x] 4.4 Implement image reference extraction from Docker container metadata (image name + tag parsing)
- [x] 4.5 Write unit tests for image reference extraction (tagged, untagged, with registry)

## 5. Worker Pool

- [x] 5.1 Implement `internal/worker/pool.go` with Pool struct (scanCh, store, trivyClient, activeWorkers atomic, logger)
- [x] 5.2 Implement Run(ctx) — N goroutines in errgroup reading from scanCh, calling SetScanning/Trivy/UpdateResult or SetError
- [x] 5.3 Track active worker count via atomic int for the vordr_scan_workers_active metric
- [x] 5.4 Write unit tests for pool (mock trivy client, verify store state transitions)

## 6. Scheduler

- [x] 6.1 Implement `internal/scheduler/scheduler.go` with Scheduler struct (interval, store, scanCh, logger)
- [x] 6.2 Implement Run(ctx) — ticker at ScanInterval, on tick: store.Images() → send all to scanCh
- [x] 6.3 Write unit test for scheduler (verify tick triggers rescan of known images)

## 7. Prometheus Metrics

- [x] 7.1 Rewrite `internal/metrics/metrics.go` — implement prometheus.Collector interface (Describe + Collect)
- [x] 7.2 In Collect(): call Store.Snapshot(), generate per-image metrics (vulnerability_count, fixable_count, scan_status, scan_timestamp, scan_duration, container_image_info)
- [x] 7.3 In Collect(): generate operational summary metrics (images_discovered, images_pending, images_scanned, images_errors)
- [x] 7.4 Add scan_workers_active and scan_queue_length gauges (fed via atomic counters from pool/channel)
- [x] 7.5 Write unit tests for collector (mock store snapshot, verify metric output)

## 8. Main Integration

- [x] 8.1 Update `cmd/vordr/main.go` — create Docker client, store, scanCh, all components
- [x] 8.2 Wire errgroup lifecycle: metrics server, watcher, scheduler, worker pool
- [x] 8.3 Pass scan_workers_active and scan_queue_length sources to metrics collector
- [x] 8.4 Verify `go build ./...` and `go vet ./...` pass

## 9. Alloy Integration

- [x] 9.1 Create Alloy scrape config snippet for vordr metrics endpoint (documentation/example file)

## 10. Verification

- [x] 10.1 Run `go test ./...` — all tests pass
- [x] 10.2 Run `go vet ./...` — no issues
- [ ] 10.3 Manual smoke test: `docker compose up`, verify /metrics returns vordr_* metrics
