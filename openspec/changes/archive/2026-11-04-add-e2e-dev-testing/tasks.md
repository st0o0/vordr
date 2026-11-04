## 1. Dev Stack

- [x] 1.1 Create `docker-compose.dev.yml` with vordr service (build: ., port 9550, docker.sock ro, trivy-cache volume, dev env vars)
- [x] 1.2 Add workload-clean (alpine:3.21), workload-vuln (nginx:1.24.0), workload-exclude (redis:7) services
- [x] 1.3 Add `E2E-TEST-RESULTS-*.md` to `.gitignore`

## 2. Test Plan

- [x] 2.1 Create `E2E-TEST-PLAN.md` with prerequisites and setup instructions
- [x] 2.2 Write Section 0: Setup verification (healthz, log output)
- [x] 2.3 Write Section 1: Initial scan (metrics for alpine/nginx, vuln counts)
- [x] 2.4 Write Section 2: Container lifecycle (stop/start tracking)
- [x] 2.5 Write Section 3: Exclude logic (redis never in metrics)
- [x] 2.6 Write Section 4: Rescan cycle (timestamp updates after interval)
- [x] 2.7 Write Section 5: Error handling (invalid image ref)
- [x] 2.8 Write Section 6: Healthcheck (HTTP and Docker inspect)
- [x] 2.9 Write Section 7: Metrics shape (vordr_ names, labels)
- [x] 2.10 Write Section 8: Details endpoint (all images, filter by image)

## 3. E2E Verify Skill

- [x] 3.1 Create `.claude/skills/e2e-verify/skill.md` with skill description and trigger
- [x] 3.2 Write Phase 0: stack lifecycle (down -v, up --build, healthz poll, scan completion poll)
- [x] 3.3 Write test execution instructions referencing E2E-TEST-PLAN.md sections
- [x] 3.4 Write results file format and naming convention (E2E-TEST-RESULTS-YYYY-MM-DD.md)

## 4. Verification

- [x] 4.1 Run `docker compose -f docker-compose.dev.yml config` to validate compose file
- [x] 4.2 Run `docker compose -f docker-compose.dev.yml up -d --build` and verify all containers start
- [x] 4.3 Execute `/e2e-verify` skill and confirm results file is produced
