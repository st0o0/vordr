## 1. Config Changes

- [x] 1.1 Replace `TrivyAddr string` with `TrivyCacheDir string` and `TrivyTimeout time.Duration` in Config struct
- [x] 1.2 Replace `VORDR_TRIVY_ADDR` parsing with `VORDR_TRIVY_CACHE_DIR` (default `/var/cache/trivy`) and `VORDR_TRIVY_TIMEOUT` (default `5m`)
- [x] 1.3 Remove net import if no longer needed after TrivyAddr removal

## 2. Trivy CLI Client

- [x] 2.1 Rewrite `internal/trivy/client.go` — replace HTTP Client with CLI struct holding cacheDir, timeout
- [x] 2.2 Implement Scan(ctx, imageRef) via exec.CommandContext: `trivy image --format json --scanners vuln --quiet --cache-dir <dir> <imageRef>`
- [x] 2.3 Parse JSON stdout into ScanResult (reuse existing Counts/Fixable maps, same severity extraction logic)
- [x] 2.4 Handle errors: non-zero exit code (include stderr), binary not found, context timeout
- [x] 2.5 Rewrite `internal/trivy/client_test.go` — test JSON parsing with sample Trivy output, test error cases

## 3. Main Integration

- [x] 3.1 Update `cmd/vordr/main.go` — replace `trivy.NewClient(cfg.TrivyAddr)` with `trivy.NewCLI(cfg.TrivyCacheDir, cfg.TrivyTimeout)`
- [x] 3.2 Remove Docker client import if no longer used for Trivy (still needed for Watcher)

## 4. Dockerfile

- [x] 4.1 Add `FROM aquasec/trivy:0.76.1 AS trivy` stage to extract Trivy binary
- [x] 4.2 Change final stage from `scratch` to `alpine:3.21`
- [x] 4.3 Add `RUN apk add --no-cache ca-certificates` to final stage
- [x] 4.4 Add `COPY --from=trivy /usr/local/bin/trivy /usr/local/bin/trivy` to final stage
- [x] 4.5 Update OCI labels if needed, keep healthcheck

## 5. Docker Compose

- [x] 5.1 Remove trivy service from docker-compose.yml
- [x] 5.2 Remove `depends_on: trivy` from vordr service
- [x] 5.3 Replace `VORDR_TRIVY_ADDR` env with `VORDR_TRIVY_CACHE_DIR: /var/cache/trivy`
- [x] 5.4 Add `trivy-cache:/var/cache/trivy` volume mount to vordr service

## 6. Documentation Updates

- [x] 6.1 Update CLAUDE.md architecture section (single container, no Trivy sidecar)
- [x] 6.2 Update ROADMAP.md configuration section (new env vars)

## 7. Verification

- [x] 7.1 Run `go build ./...` and `go vet ./...` — pass
- [x] 7.2 Run `go test ./...` — all tests pass
- [x] 7.3 Run `docker compose build` — image builds successfully
- [x] 7.4 Run `docker compose up`, verify /metrics returns vordr_* metrics with real scan data
