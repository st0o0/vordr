## ADDED Requirements

### Requirement: Environment variable based configuration
The config package SHALL load all configuration from environment variables prefixed with `VORDR_`. The `Load` function SHALL accept a `getenv func(string) string` parameter for testability.

#### Scenario: Default values
- **WHEN** no environment variables are set
- **THEN** Load SHALL return a Config with these defaults:
  - TrivyAddr: `http://trivy:4954`
  - ScanInterval: `6h`
  - ScanWorkers: `3`
  - ScanStagger: `30s`
  - ListenAddr: `:9550`
  - LogLevel: `info`
  - LogFormat: `text`

#### Scenario: Custom values
- **WHEN** VORDR_TRIVY_ADDR is set to `http://localhost:4954`
- **THEN** Load SHALL return a Config with TrivyAddr `http://localhost:4954`

#### Scenario: Invalid duration
- **WHEN** VORDR_SCAN_INTERVAL is set to `not-a-duration`
- **THEN** Load SHALL return a descriptive error

#### Scenario: Invalid integer
- **WHEN** VORDR_SCAN_WORKERS is set to `abc`
- **THEN** Load SHALL return a descriptive error

### Requirement: Structured logging setup
The main function SHALL create a `log/slog` logger with JSON or text format based on VORDR_LOG_FORMAT and configurable level via VORDR_LOG_LEVEL.

#### Scenario: JSON logging
- **WHEN** VORDR_LOG_FORMAT is `json`
- **THEN** log output SHALL be structured JSON

#### Scenario: Text logging
- **WHEN** VORDR_LOG_FORMAT is `text` or unset
- **THEN** log output SHALL be human-readable text
