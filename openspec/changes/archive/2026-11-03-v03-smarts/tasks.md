## 1. Config Changes

- [ ] 1.1 Add `ExcludeImages []string`, `AlertSeverity string`, `RetryMaxBackoff time.Duration` to Config struct
- [ ] 1.2 Parse `VORDR_EXCLUDE_IMAGES` (comma-split, trim spaces), `VORDR_ALERT_SEVERITY` (validate against known severities, default `HIGH`), `VORDR_RETRY_MAX_BACKOFF` (default `15m`)
- [ ] 1.3 Add `ShouldExclude(imageTag string) bool` method on Config using path.Match

## 2. Image Exclude Patterns

- [ ] 2.1 Update `internal/watcher/watcher.go` Run() — skip excluded images during initial listing, log `event=image_excluded`
- [ ] 2.2 Update `internal/watcher/watcher.go` handleEvent() — skip excluded images on start events
- [ ] 2.3 Update `internal/scheduler/scheduler.go` rescan() — skip excluded images when queuing
- [ ] 2.4 Pass Config (or exclude func) to Watcher and Scheduler constructors
- [ ] 2.5 Write tests for ShouldExclude (matching, non-matching, empty patterns, multiple patterns)

## 3. Delta Detection

- [ ] 3.1 Add `PreviousVulnIDs map[string]bool` to store entry, save current IDs as previous before overwrite in UpdateResult
- [ ] 3.2 Add `PreviousVulnIDs map[string]bool` to `scanner.ImageState` and populate in Snapshot()
- [ ] 3.3 Add delta detection logic in worker: after scan, compute new CVE IDs = current minus previous
- [ ] 3.4 Add severity ordering helper (`internal/scanner/severity.go`) with `MeetsThreshold(severity, threshold string) bool`
- [ ] 3.5 Log `event=new_cves_detected` with new_count, alert_worthy_count, new CVE IDs when delta is non-empty
- [ ] 3.6 Write tests for delta detection (first scan, rescan with new CVEs, rescan with same CVEs)
- [ ] 3.7 Write tests for severity threshold

## 4. Retry Backoff

- [ ] 4.1 Create `internal/retry/queue.go` — RetryQueue struct with backoff logic (30s base, 2x multiplier, max from config)
- [ ] 4.2 Implement `Add(ref ImageRef)` — schedules retry with current backoff, `Reset(image string)` — clears backoff on success
- [ ] 4.3 Implement `Run(ctx) error` — goroutine that waits for backoff then sends to scanCh
- [ ] 4.4 Wire retry queue: worker calls `retryQueue.Add(ref)` on scan failure, `retryQueue.Reset(key)` on success
- [ ] 4.5 Wire retry queue into errgroup in main.go
- [ ] 4.6 Write tests for retry queue (backoff calculation, max cap, reset)

## 5. Main Integration

- [ ] 5.1 Update `cmd/vordr/main.go` — create retry queue, pass exclude func and retry queue to components
- [ ] 5.2 Update startup log to include new config values

## 6. Documentation

- [ ] 6.1 Update ROADMAP.md configuration section with new env vars
- [ ] 6.2 Update CLAUDE.md if needed

## 7. Verification

- [ ] 7.1 Run `go build ./...` and `go vet ./...`
- [ ] 7.2 Run `go test ./...` — all tests pass
- [ ] 7.3 Docker compose build + smoke test
