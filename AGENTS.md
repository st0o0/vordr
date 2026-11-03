# Agents

## Code changes

Read `CLAUDE.md` for project conventions before making changes. Run `go build ./...` and `go test ./...` after edits. Run `golangci-lint run` before committing.

## OpenSpec

Use `/opsx:explore` for design discussions. Use `/opsx:propose` to create changes. Use `/opsx:apply` to implement. Config in `openspec/config.yaml`.

## Security-sensitive areas

- `internal/trivy/` -- HTTP client communicating with Trivy server
- `internal/watcher/` -- Docker socket access (read-only)
- `internal/store/` -- shared mutable state (sync.RWMutex)
