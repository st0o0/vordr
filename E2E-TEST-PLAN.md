# E2E Test Plan

Complete API and metrics test plan for vordr running in a dev Docker stack.
Do not check items off in this file. Write results to `E2E-TEST-RESULTS-<DATE>.md`.

## Prerequisites

- Docker Desktop or Docker Engine with compose plugin
- Port 9550 available
- Internet access (Trivy DB download on first run)

## Stack

```bash
docker compose -f docker-compose.dev.yml up -d --build
```

| Service | Image | Container | Purpose |
|---------|-------|-----------|---------|
| vordr | build from source | vordr | System under test |
| workload-clean | alpine:3.21 | workload-clean | CVE-free baseline |
| workload-vuln | nginx:1.24.0 | workload-vuln | Known CVEs |
| workload-exclude | redis:7 | workload-exclude | Excluded via VORDR_EXCLUDE_IMAGES=redis* |

### Dev Environment Variables

| Variable | Value | Purpose |
|----------|-------|---------|
| VORDR_SCAN_INTERVAL | 30s | Fast rescan cycle for testing |
| VORDR_SCAN_STAGGER | 2s | Short delay between initial scans |
| VORDR_SCAN_WORKERS | 2 | Parallel scan workers |
| VORDR_LOG_LEVEL | debug | Verbose output for debugging |
| VORDR_LOG_FORMAT | text | Human-readable logs |
| VORDR_EXCLUDE_IMAGES | redis*,vordr-* | Excludes the redis workload and vordr's own image |

### Endpoints

| Path | Description |
|------|-------------|
| `GET http://localhost:9550/healthz` | Health check (JSON: status, version, uptime) |
| `GET http://localhost:9550/metrics` | Prometheus metrics |
| `GET http://localhost:9550/details` | All scanned images with vulnerability details (JSON) |
| `GET http://localhost:9550/details?image=<ref>` | Filtered by image reference |

### Timing Notes

- First run: Trivy downloads its vulnerability DB (~40MB). Allow up to 3 minutes.
- Subsequent runs with `trivy-cache` volume intact: scans start immediately.
- `VORDR_SCAN_INTERVAL=30s`: rescan cycle triggers every 30 seconds.
- Use `docker compose -f docker-compose.dev.yml down` (without `-v`) to preserve the cache.
- Use `docker compose -f docker-compose.dev.yml down -v` for a full reset including cache.

---

## Section 0: Setup Verification

Verify the dev stack is operational before running functional tests.

### 0.1 All containers running

```bash
docker compose -f docker-compose.dev.yml ps --format json
```

- [ ] vordr container is running and healthy
- [ ] workload-clean container is running
- [ ] workload-vuln container is running
- [ ] workload-exclude container is running

### 0.2 Health endpoint responds

```
GET http://localhost:9550/healthz
```

- [ ] Returns HTTP 200
- [ ] JSON body contains `"status": "ok"`
- [ ] JSON body contains `"version"` field

### 0.3 Logs indicate startup

```bash
docker compose -f docker-compose.dev.yml logs vordr --tail 50
```

- [ ] Logs contain structured entries indicating vordr started
- [ ] Logs show scan workers starting
- [ ] No ERROR-level entries during startup

---

## Section 1: Initial Scan

Wait for vordr to complete initial scans of discovered containers. Poll `/metrics` until `vordr_images_scanned` shows at least 2 images (alpine + nginx). Timeout: 5 minutes.

### 1.1 Non-excluded images in metrics

```
GET http://localhost:9550/metrics
```

- [ ] Metrics contain `vordr_image_scan_status` for alpine image
- [ ] Metrics contain `vordr_image_scan_status` for nginx image
- [ ] `vordr_images_discovered` is 2 (only non-excluded images)
- [ ] `vordr_images_scanned` is 2

### 1.2 Vulnerable image has CVEs

```
GET http://localhost:9550/metrics
```

- [ ] `vordr_vulnerability_count{image="nginx",tag="1.24.0",severity=...}` has at least one severity with count > 0

### 1.3 Clean image has zero or minimal CVEs

```
GET http://localhost:9550/metrics
```

- [ ] `vordr_vulnerability_count` for alpine:3.21 shows zero or very low counts across all severities

### 1.4 Container-to-image mapping

```
GET http://localhost:9550/metrics
```

- [ ] `vordr_container_image_info{container="workload-clean",image="alpine",tag="3.21"}` exists
- [ ] `vordr_container_image_info{container="workload-vuln",image="nginx",tag="1.24.0"}` exists

### 1.5 Scan metadata present

```
GET http://localhost:9550/metrics
```

- [ ] `vordr_scan_timestamp_seconds` is present for both scanned images
- [ ] `vordr_scan_duration_seconds` is present for both scanned images and values are > 0

---

## Section 2: Container Lifecycle

Test that vordr reacts to container stop and start events via Docker API.

### 2.1 Stop removes image from tracking

```bash
docker compose -f docker-compose.dev.yml stop workload-vuln
```

Wait 5 seconds, then:

```
GET http://localhost:9550/metrics
```

- [ ] `vordr_image_scan_status` for nginx is no longer present
- [ ] `vordr_container_image_info` for workload-vuln is no longer present
- [ ] `vordr_images_discovered` decreased

### 2.2 Start triggers rescan

```bash
docker compose -f docker-compose.dev.yml start workload-vuln
```

Poll `/metrics` until nginx reappears (timeout: 2 minutes).

- [ ] `vordr_image_scan_status` for nginx reappears
- [ ] `vordr_vulnerability_count` for nginx shows CVEs again
- [ ] `vordr_container_image_info` for workload-vuln is present

---

## Section 3: Exclude Logic

Verify that excluded images (matching `VORDR_EXCLUDE_IMAGES=redis*,vordr-*`) are never tracked.

### 3.1 Excluded images absent from metrics

```
GET http://localhost:9550/metrics
```

- [ ] No `vordr_image_scan_status` metric exists with `image="redis"`
- [ ] No `vordr_image_scan_status` metric exists with `image="vordr-vordr"`
- [ ] No `vordr_container_image_info` metric exists with `container="workload-exclude"`

### 3.2 Excluded images absent from details

```
GET http://localhost:9550/details
```

- [ ] JSON response does not contain any entry with `"image": "redis"`
- [ ] JSON response does not contain any entry with `"image": "vordr-vordr"`

---

## Section 4: Rescan Cycle

Verify that periodic rescanning updates scan timestamps.

### 4.1 Record initial scan timestamp

```
GET http://localhost:9550/metrics
```

- [ ] Record `vordr_scan_timestamp_seconds` for alpine

### 4.2 Wait for rescan and verify timestamp updated

Wait at least 35 seconds (VORDR_SCAN_INTERVAL=30s + margin), then:

```
GET http://localhost:9550/metrics
```

- [ ] `vordr_scan_timestamp_seconds` for alpine is greater than the value recorded in 4.1

---

## Section 5: Error Handling

Test that vordr handles scan failures gracefully.

### 5.1 Start container with invalid image

```bash
docker run -d --name workload-broken --label vordr.test=true alpine:3.21 sh -c "sleep infinity"
docker rename workload-broken workload-error-test
```

Note: vordr tracks by image reference. To test error handling, we can observe the error state via logs when a scan fails. Check `docker compose -f docker-compose.dev.yml logs vordr` for any error-level entries related to scan failures.

Alternatively, check the details endpoint for any images in error state:

```
GET http://localhost:9550/details
```

- [ ] If any image is in error state, `"status": "error"` is visible in the details response
- [ ] Error does not crash vordr (healthz still returns 200)

### 5.2 Cleanup

```bash
docker rm -f workload-error-test 2>/dev/null
```

---

## Section 6: Healthcheck

### 6.1 HTTP health endpoint

```
GET http://localhost:9550/healthz
```

- [ ] Returns HTTP 200
- [ ] JSON body has `"status": "ok"`
- [ ] `"uptime"` value is a non-zero duration string

### 6.2 Docker healthcheck status

```bash
docker inspect vordr --format "{{.State.Health.Status}}"
```

- [ ] Returns `healthy`

---

## Section 7: Metrics Shape

Verify all expected `vordr_` metric families are present with correct labels.

### 7.1 All metric families present

```
GET http://localhost:9550/metrics
```

Expected metric families (all prefixed `vordr_`):

- [ ] `vordr_image_scan_status` (labels: image, tag)
- [ ] `vordr_vulnerability_count` (labels: image, tag, severity)
- [ ] `vordr_fixable_count` (labels: image, tag, severity)
- [ ] `vordr_scan_timestamp_seconds` (labels: image, tag)
- [ ] `vordr_scan_duration_seconds` (labels: image, tag)
- [ ] `vordr_container_image_info` (labels: container, image, tag)
- [ ] `vordr_images_discovered` (no labels)
- [ ] `vordr_images_pending` (no labels)
- [ ] `vordr_images_scanned` (no labels)
- [ ] `vordr_images_errors` (no labels)
- [ ] `vordr_scan_workers_active` (no labels)
- [ ] `vordr_scan_queue_length` (no labels)
- [ ] `vordr_trivy_db_updated_timestamp_seconds` (no labels)

### 7.2 Severity labels are valid

```
GET http://localhost:9550/metrics
```

- [ ] `severity` label values are from the set: UNKNOWN, LOW, MEDIUM, HIGH, CRITICAL

---

## Section 8: Details Endpoint

### 8.1 All scanned images returned

```
GET http://localhost:9550/details
```

- [ ] Response is a JSON array
- [ ] Contains entry for alpine with `"status": "scanned"`
- [ ] Contains entry for nginx with `"status": "scanned"`
- [ ] Does not contain entry for redis

### 8.2 Vulnerability details for nginx

```
GET http://localhost:9550/details
```

- [ ] nginx entry has non-empty `"vulnerabilities"` array
- [ ] Each vulnerability has fields: `id`, `package`, `installed`, `severity`
- [ ] At least one vulnerability has a non-empty `"fixed"` field

### 8.3 Filter by image

```
GET http://localhost:9550/details?image=nginx:1.24.0
```

- [ ] Response contains only the nginx entry
- [ ] Returns HTTP 200

### 8.4 Filter with unknown image returns 404

```
GET http://localhost:9550/details?image=nonexistent:latest
```

- [ ] Returns HTTP 404
- [ ] JSON body contains `"error": "image not found"`

---

## Test Execution Order

Execute sections in order 0 through 8. Each section builds on the state established by previous sections.

1. Section 0: Setup verification (must pass before continuing)
2. Section 1: Initial scan (wait for scans to complete)
3. Section 2: Container lifecycle (stop/start workload-vuln)
4. Section 3: Exclude logic (verify redis exclusion)
5. Section 4: Rescan cycle (requires ~35s wait)
6. Section 5: Error handling
7. Section 6: Healthcheck
8. Section 7: Metrics shape
9. Section 8: Details endpoint

## Teardown

```bash
docker compose -f docker-compose.dev.yml down -v
```
