## Context

vordr has unit tests for individual components but no way to verify the system works end-to-end. The critical paths - Docker event handling, Trivy CLI scanning, metrics emission - cross component boundaries and are only testable with a real Docker environment. The existing `docker-compose.yml` is production-oriented (6h scan interval, no test workloads). We need a dev-focused compose file with fast feedback loops and a structured test plan.

The funkarr project established a pattern: a reusable `E2E-TEST-PLAN.md` as the spec, and a `.claude/skills/e2e-verify/` skill as the orchestrator that executes the plan and writes results. This pattern works well for vordr, adapted to its simpler API-only surface.

## Goals / Non-Goals

**Goals:**
- Verify vordr's observable behavior end-to-end with real Docker and Trivy
- Reproducible dev stack that starts with one command
- Structured test plan covering all HTTP endpoints and container lifecycle events
- Skill-driven execution producing consistent, comparable results files

**Non-Goals:**
- CI integration (local-only for now)
- Performance or load testing
- Testing Prometheus scraping or alerting (Level 3 scope)
- Replacing existing unit tests

## Decisions

### D1: Three test workload containers

**Decision:** Use `alpine:3.21` (clean), `nginx:1.24.0` (vulnerable), `redis:7` (excluded).

**Rationale:** These cover the three key states vordr must handle: a scanned image with no CVEs, a scanned image with known CVEs, and an image that matches `VORDR_EXCLUDE_IMAGES` and must be ignored. Using well-known public images ensures reproducibility without maintaining custom test images.

**Alternative considered:** Building custom test images with planted CVEs. Rejected - unnecessary complexity, and nginx:1.24.0 reliably has known CVEs.

### D2: Build from source in compose

**Decision:** `build: .` for the vordr service, using the existing multi-stage Dockerfile.

**Rationale:** Always tests the current code. The Dockerfile already handles bundling the Trivy binary, so no separate Trivy installation needed in the dev stack.

**Alternative considered:** Pre-built image from GHCR. Rejected - defeats the purpose of testing local changes.

### D3: Aggressive timing for dev feedback

**Decision:** `VORDR_SCAN_INTERVAL=30s`, `VORDR_SCAN_STAGGER=2s`, `VORDR_SCAN_WORKERS=2`.

**Rationale:** Production defaults (6h interval, 30s stagger) are too slow for interactive testing. 30s interval means rescan behavior is observable within a minute. 2s stagger keeps initial scans fast.

### D4: Sequential test execution, no subagents

**Decision:** The e2e-verify skill runs all tests sequentially from the main agent.

**Rationale:** vordr has ~8 test areas with no complex state dependencies between them. The tests are all curl/Invoke-RestMethod based (no browser). Subagent waves add orchestration complexity that isn't justified for this scope. Unlike funkarr's 43-step plan with browser + API + SSE tests, vordr's plan is simple enough for linear execution.

### D5: Skill reads plan as reusable spec

**Decision:** `E2E-TEST-PLAN.md` is the source of truth for test definitions. The skill references it but never modifies it. Results go to `E2E-TEST-RESULTS-<DATE>.md`.

**Rationale:** Follows the funkarr pattern. The plan is versioned in git and evolves with the codebase. Results files are per-run and gitignored. This separation means the plan can be reviewed independently of execution.

### D6: Poll-based scan completion detection

**Decision:** After compose up, poll `/metrics` until both non-excluded images appear with scan results, with a timeout accommodating Trivy DB download.

**Rationale:** The initial scan involves Trivy downloading its vulnerability DB on first run (~40MB, 1-2 min) plus the actual scans. A fixed sleep would either be too short or wastefully long. Polling `/metrics` for the expected image entries is deterministic.

**Timeout:** 5 minutes for initial scan completion (covers DB download + scan time). Subsequent scans should complete within the 30s interval.

## Risks / Trade-offs

**[Trivy DB download variability]** → First run depends on network speed for DB download. Mitigated by trivy-cache named volume that persists across `compose down` (only `down -v` clears it). The skill uses `-v` for a clean slate but documents this tradeoff.

**[nginx:1.24.0 CVE count may change]** → Trivy DB updates may add/remove CVEs for this image. Tests should assert `count > 0` rather than exact numbers.

**[Docker Desktop vs Docker Engine differences]** → The containerd image store bug (see memory) affects digest-pulled images. Our test workloads use tags, not digests, so this should not apply. Test plan should note this as a known limitation.

**[Port 9550 conflict]** → If the user runs vordr locally and in compose simultaneously. Mitigated by using the same port as production compose, making the conflict obvious.
