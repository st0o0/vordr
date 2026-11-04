## Context

vordr currently runs Trivy as a CLI subprocess with a local vulnerability database (3.1 GB). In multi-machine deployments this means every instance downloads and stores its own copy. Trivy supports a client-server mode where the server holds the DB and clients send only package metadata for lookup. The JSON output format is identical, making the switch transparent to the parsing layer.

The CI uses shared workflows from `st0o0/github-workflows`: `docker-release.yml` for production builds and `docker-dev.yml` for PR dev builds. Both accept `image-name` and `image-description` parameters.

## Goals / Non-Goals

**Goals:**
- Support both standalone and server mode, controlled by a single env var
- Ship a companion `vordr-server` image for easy deployment
- Build and publish both images in CI (release + dev-build)
- Backwards compatible — no config change = identical behavior to today
- Reduce standalone DB size with `--skip-java-db-update` (1.5 GB savings)

**Non-Goals:**
- Trivy server authentication (token support) — can be added later if needed
- Automatic fallback from server to standalone mode
- Health-checking the Trivy server from vordr (scans will simply fail, retry queue handles it)
- Changing the Trivy version pinning strategy

## Decisions

### D1: Single env var controls mode

**Decision:** `VORDR_TRIVY_SERVER_URL` — empty or unset = standalone, set = server mode.

**Rationale:** Follows the existing pattern where each VORDR_ env var controls one behavior. No boolean flag needed — the URL presence is the signal. This avoids the ambiguity of having both a server URL and a cache dir set.

**Alternative considered:** Two vars (VORDR_TRIVY_MODE + VORDR_TRIVY_SERVER_URL). Rejected — the URL alone is sufficient and less error-prone.

### D2: CLI constructor accepts optional server URL

**Decision:** Change `trivy.NewCLI(cacheDir, timeout)` to `trivy.NewCLI(cacheDir, timeout, serverURL)` where serverURL is empty string for standalone mode.

**Rationale:** Keeps the mode decision inside the trivy package. The CLI struct holds both cacheDir and serverURL; Scan() picks args based on which is set. The cacheDir field is still needed in standalone mode.

**Alternative considered:** Two separate types (StandaloneCLI, ServerCLI) behind the Scanner interface. Rejected — the only difference is 2 CLI args; a type split is overengineering.

### D3: DBChecker receives a disabled flag

**Decision:** `trivy.NewDBChecker(cacheDir, interval, logger, enabled)` — in server mode, pass `enabled=false`. The Run() method returns immediately when disabled.

**Rationale:** The DBChecker's `trivy version --format json` command reads the local DB. In server mode there is no local DB, so the command would fail. The metric `vordr_trivy_db_updated_timestamp_seconds` simply won't be emitted — this is a correct representation (vordr doesn't know the server's DB age).

### D4: Dockerfile.server is a thin wrapper

**Decision:** Single-stage Dockerfile based on `aquasec/trivy:0.76.1` with `ENTRYPOINT ["trivy", "server", "--listen", "0.0.0.0:4954"]`. Includes a HEALTHCHECK.

**Rationale:** The server image is just Trivy configured to run in server mode. No custom code needed. The version pin matches the trivy stage in the main Dockerfile — one place to update.

### D5: CI adds a second docker job for vordr-server

**Decision:** Both `release.yml` and `dev-build.yml` get a second job (`docker-server`) that calls the same shared workflow with `image-name: ghcr.io/st0o0/vordr-server` and `dockerfile: Dockerfile.server`.

**Rationale:** Reuses the existing shared workflow. Both images get the same version tag from release-please. The server job runs in parallel with the vordr job since they're independent.

### D6: --skip-java-db-update in standalone mode

**Decision:** Add `--skip-java-db-update` to standalone Scan() args. Not applicable in server mode (server controls its own DB flags).

**Rationale:** Saves 1.5 GB per standalone instance. Most homelab stacks don't run Java containers. If a user needs Java scanning, they can switch to server mode where the server operator controls DB flags.

**Trade-off:** Java-specific CVEs may be missed in standalone mode. The main vuln DB still catches OS-level and many library CVEs for Java images.

## Risks / Trade-offs

**[Version mismatch]** → Trivy client and server must run the same version. Mitigated by pinning both to the same version in their respective Dockerfiles, and releasing both images together.

**[Single point of failure]** → Trivy server down means all scans fail. Mitigated by vordr's existing retry queue — scans will retry with backoff until the server recovers.

**[DB metric missing in server mode]** → `vordr_trivy_db_updated_timestamp_seconds` won't be emitted. This is correct — vordr doesn't know the server's DB age. Users monitoring this metric need to check the Trivy server directly.

**[Shared workflow may not support dockerfile parameter]** → The `docker-release.yml` and `docker-dev.yml` shared workflows may need a `dockerfile` input. If not supported, we'll need to add it to the shared workflows or use a `context` parameter. Check before implementing.
