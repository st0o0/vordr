## ADDED Requirements

### Requirement: Go module initialization
The project SHALL have a Go module at `github.com/st0o0/vordr` with Go 1.26.x and dependencies on docker/docker, prometheus/client_golang, and golang.org/x/sync.

#### Scenario: Valid go.mod
- **WHEN** the repository is cloned
- **THEN** `go mod download` SHALL succeed without errors

#### Scenario: Module path
- **WHEN** inspecting go.mod
- **THEN** the module path SHALL be `github.com/st0o0/vordr`

### Requirement: Binary entry point with subcommands
The project SHALL have a `cmd/vordr/main.go` that builds a single binary supporting `version` and `healthcheck` subcommands. The binary SHALL inject version via `-ldflags "-X main.version=${VERSION}"`.

#### Scenario: Version subcommand
- **WHEN** running `vordr version`
- **THEN** the binary SHALL print the version string and exit 0

#### Scenario: Healthcheck subcommand
- **WHEN** running `vordr healthcheck` and the metrics server is reachable
- **THEN** the binary SHALL HTTP GET `/healthz` on the configured metrics address and exit 0 on 200, exit 1 otherwise

#### Scenario: Default run
- **WHEN** running `vordr` without subcommands
- **THEN** the binary SHALL start the application (config loading, logger setup, metrics server)

### Requirement: Internal package structure
The project SHALL have internal packages that define the architecture: config/, metrics/, watcher/, scheduler/, worker/, store/, trivy/.

#### Scenario: Package compilation
- **WHEN** running `go build ./...`
- **THEN** all packages SHALL compile without errors

#### Scenario: Stub packages
- **WHEN** inspecting watcher/, scheduler/, worker/, store/, trivy/
- **THEN** each SHALL contain at minimum a Go file with only the package declaration

### Requirement: Signal handling and graceful shutdown
The main function SHALL use `signal.NotifyContext` with SIGTERM and SIGINT for graceful shutdown.

#### Scenario: Clean shutdown
- **WHEN** the process receives SIGTERM
- **THEN** the context SHALL be cancelled and all goroutines SHALL return cleanly
