## MODIFIED Requirements

### Requirement: Environment variable based configuration
The config package SHALL load all configuration from environment variables prefixed with `VORDR_`. The `Load` function SHALL accept a `getenv func(string) string` parameter for testability.

#### Scenario: Default values
- **WHEN** no environment variables are set
- **THEN** Load SHALL return a Config with these defaults:
  - TrivyCacheDir: `/var/cache/trivy`
  - TrivyTimeout: `5m`
  - ScanInterval: `6h`
  - ScanWorkers: `3`
  - ScanStagger: `30s`
  - ListenAddr: `:9550`
  - LogLevel: `info`
  - LogFormat: `text`

#### Scenario: Custom cache directory
- **WHEN** VORDR_TRIVY_CACHE_DIR is set to `/data/trivy`
- **THEN** Load SHALL return a Config with TrivyCacheDir `/data/trivy`

#### Scenario: Custom timeout
- **WHEN** VORDR_TRIVY_TIMEOUT is set to `10m`
- **THEN** Load SHALL return a Config with TrivyTimeout 10 minutes

#### Scenario: Invalid timeout
- **WHEN** VORDR_TRIVY_TIMEOUT is set to `not-a-duration`
- **THEN** Load SHALL return a descriptive error

#### Scenario: Removed VORDR_TRIVY_ADDR
- **WHEN** VORDR_TRIVY_ADDR is set
- **THEN** it SHALL be ignored (no longer used)
