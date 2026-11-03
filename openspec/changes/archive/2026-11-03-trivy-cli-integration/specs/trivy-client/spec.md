## MODIFIED Requirements

### Requirement: Trivy scan execution
The Trivy client SHALL execute `trivy image --format json --scanners vuln --quiet --cache-dir <dir> <imageRef>` as a subprocess via `exec.CommandContext` and parse the JSON output to extract vulnerability counts by severity.

#### Scenario: Successful scan
- **WHEN** scanning `nginx:latest` and Trivy exits 0 with JSON output
- **THEN** the client SHALL return a ScanResult with Counts and Fixable maps populated from Results[].Vulnerabilities[]

#### Scenario: Image with no vulnerabilities
- **WHEN** scanning an image and Trivy returns empty Results or no Vulnerabilities
- **THEN** the client SHALL return zero counts for all severities

#### Scenario: Trivy binary not found
- **WHEN** the trivy binary is not on PATH
- **THEN** the client SHALL return an error

#### Scenario: Trivy exits with non-zero code
- **WHEN** Trivy exits with a non-zero exit code
- **THEN** the client SHALL return an error including stderr output

#### Scenario: Scan timeout
- **WHEN** a scan exceeds the configured timeout
- **THEN** the context SHALL cancel the subprocess and the client SHALL return a timeout error

### Requirement: HTTP client configuration
The Trivy client SHALL use the configured cache directory and timeout. No HTTP connection to a Trivy server is needed.

#### Scenario: Custom cache directory
- **WHEN** VORDR_TRIVY_CACHE_DIR is `/data/trivy`
- **THEN** the client SHALL pass `--cache-dir /data/trivy` to the trivy command

#### Scenario: Request timeout
- **WHEN** VORDR_TRIVY_TIMEOUT is 10m
- **THEN** the context timeout for exec.CommandContext SHALL be 10 minutes
