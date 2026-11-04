# trivy-server-mode Specification

## Purpose
TBD - created by archiving change add-trivy-server-mode. Update Purpose after archive.
## Requirements
### Requirement: Dual-mode Trivy scanning
vordr SHALL support two scanning modes: standalone (local Trivy DB) and server (remote Trivy server for DB lookups). The mode SHALL be determined by the presence of the VORDR_TRIVY_SERVER_URL environment variable.

#### Scenario: Standalone mode (default)
- **WHEN** VORDR_TRIVY_SERVER_URL is not set or empty
- **THEN** vordr uses the local Trivy binary with `--cache-dir` and a local vulnerability database

#### Scenario: Server mode
- **WHEN** VORDR_TRIVY_SERVER_URL is set to a valid URL (e.g., `http://trivy-server:4954`)
- **THEN** vordr uses the local Trivy binary with `--server <url>` instead of `--cache-dir`
- **AND** no local vulnerability database is downloaded or maintained

#### Scenario: Scan results are identical in both modes
- **WHEN** the same image is scanned in standalone and server mode
- **THEN** the JSON output format and parsed ScanResult structure are identical

