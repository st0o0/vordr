## Why

Each vordr instance bundles a local Trivy binary and downloads a 3.1 GB vulnerability database (1.6 GB main DB + 1.5 GB Java DB). In a multi-machine homelab this scales linearly — 8 machines means ~25 GB of duplicated DB storage plus 8 independent DB update cycles. A centralized Trivy server lets all instances share one DB, reducing total storage to ~3.1 GB and eliminating redundant downloads.

## What Changes

- Add `VORDR_TRIVY_SERVER_URL` env var for dual-mode operation: empty = standalone (current behavior), set = server mode (Trivy client connects to central server)
- Modify `internal/trivy/client.go` to switch CLI args based on mode (`--cache-dir` vs `--server`)
- Modify `internal/trivy/dbcheck.go` to skip when server mode is active (server manages its own DB)
- Add `Dockerfile.server` — thin wrapper around `aquasec/trivy` with server entrypoint, published as `ghcr.io/st0o0/vordr-server`
- Add vordr-server to CI: release workflow builds and pushes both images, dev-build workflow builds both images
- Update `docker-compose.dev.yml` with vordr-server service for E2E testing of server mode
- Add `--skip-java-db-update` to standalone mode scans (saves 1.5 GB when no Java images are present)

## Capabilities

### New Capabilities
- `trivy-server-mode`: Dual-mode Trivy scanning — standalone with local DB or client mode connecting to a central Trivy server
- `server-image`: Dockerfile and CI pipeline for the vordr-server companion image

### Modified Capabilities
- `trivy-client`: Scan() switches CLI args based on mode; adds `--skip-java-db-update` for standalone
- `config-loading`: New VORDR_TRIVY_SERVER_URL env var
- `trivy-db-monitoring`: DBChecker skips in server mode (server manages its own DB)
- `build-and-ci`: Release and dev-build workflows build both vordr and vordr-server images
- `deployment`: Documentation of server mode deployment pattern

## Impact

- `internal/trivy/client.go` — constructor and Scan() args change
- `internal/trivy/dbcheck.go` — conditional skip based on server mode
- `internal/config/config.go` — new TrivyServerURL field
- `cmd/vordr/main.go` — wiring: pass server URL to CLI constructor, conditionally skip DBChecker
- `Dockerfile.server` — new file
- `.github/workflows/release.yml` — second docker job for vordr-server
- `.github/workflows/dev-build.yml` — second docker job for vordr-server
- `docker-compose.dev.yml` — vordr-server service added
- `docker-compose.yml` — document server mode option in comments
- Backwards compatible: no env var set = identical behavior to today
