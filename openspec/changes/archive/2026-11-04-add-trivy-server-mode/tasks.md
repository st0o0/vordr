## 1. Config

- [x] 1.1 Add `TrivyServerURL string` field to Config struct and load `VORDR_TRIVY_SERVER_URL` in `config.go`
- [x] 1.2 Add test case for `VORDR_TRIVY_SERVER_URL` in `config_test.go`

## 2. Trivy Client

- [x] 2.1 Add `serverURL` parameter to `trivy.NewCLI` constructor
- [x] 2.2 Modify `CLI.Scan()` to switch args: `--server <url>` when serverURL is set, `--cache-dir <dir> --skip-java-db-update` when standalone
- [x] 2.3 Add unit tests for both modes in `client_test.go` (verify command args construction)

## 3. DB Checker

- [x] 3.1 Add `enabled bool` parameter to `trivy.NewDBChecker` constructor
- [x] 3.2 Make `DBChecker.Run()` return immediately when disabled
- [x] 3.3 Add test for disabled DBChecker in `dbcheck_test.go`

## 4. Main Wiring

- [x] 4.1 Update `cmd/vordr/main.go` to pass `cfg.TrivyServerURL` to `trivy.NewCLI`
- [x] 4.2 Pass `cfg.TrivyServerURL == ""` as enabled flag to `trivy.NewDBChecker`
- [x] 4.3 Add `trivy_server_url` to startup log when set

## 5. Server Dockerfile

- [x] 5.1 Create `Dockerfile.server` based on `aquasec/trivy:0.76.1` with server entrypoint and healthcheck
- [x] 5.2 Verify `docker build -f Dockerfile.server -t vordr-server .` succeeds

## 6. CI/CD

- [x] 6.1 Add `docker-server` job to `.github/workflows/release.yml` building `ghcr.io/st0o0/vordr-server` from `Dockerfile.server`
- [x] 6.2 Add `docker-server` job to `.github/workflows/dev-build.yml` building dev `ghcr.io/st0o0/vordr-server`

## 7. Dev Stack & E2E

- [x] 7.1 Add `vordr-server` service to `docker-compose.dev.yml` with trivy-cache volume
- [x] 7.2 Add commented `VORDR_TRIVY_SERVER_URL` to vordr service in `docker-compose.dev.yml` for easy mode switching
- [x] 7.3 Add server mode comment to production `docker-compose.yml`

## 8. Verification

- [x] 8.1 Run `go test ./...` to verify all tests pass
- [x] 8.2 Run `go vet ./...` and `golangci-lint run`
- [x] 8.3 Test standalone mode: `docker compose -f docker-compose.dev.yml up -d --build` without server URL, verify scans work
- [x] 8.4 Test server mode: enable VORDR_TRIVY_SERVER_URL, rebuild, verify scans work via server
