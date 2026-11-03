# Contributing

Thanks for your interest in vordr!

## Development

vordr is a pure Go module. All packages build and test on any platform.

```bash
# tests + vet (Go)
go test ./...
go vet ./...

# lint
golangci-lint run          # golangci-lint v2
docker run --rm -i hadolint/hadolint < Dockerfile

# build
go build -o vordr ./cmd/vordr

# docker build
docker build -t vordr:ci .
```

## Pull requests

- Branch from `main` and open a PR. CI (test, lint, build) runs on pull
  requests.
- Commit messages follow [Conventional Commits](https://www.conventionalcommits.org/)
  (`feat:`, `fix:`, `docs:`, `ci:`, `chore:`, ...); commitlint enforces this and
  release-please derives the version and changelog from them.
- Keep changes focused and add or update tests for any behavior change.
- Go code must pass `go vet` and `golangci-lint`; the `Dockerfile` must pass
  `hadolint`.

## Reporting bugs / requesting features

Use the issue templates. For security issues see [`SECURITY.md`](SECURITY.md).

By contributing you agree to the [Code of Conduct](CODE_OF_CONDUCT.md).
