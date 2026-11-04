## MODIFIED Requirements

### Requirement: Trivy scan execution
The Trivy CLI client SHALL support both standalone and server mode scanning. In standalone mode, it SHALL use `--cache-dir` and `--skip-java-db-update`. In server mode, it SHALL use `--server <url>` instead. ScanResult SHALL include a Vulnerabilities slice with ID, PkgName, InstalledVersion, FixedVersion, and Severity per CVE.

#### Scenario: Successful scan with vulnerabilities
- **WHEN** scanning an image with vulnerabilities in either standalone or server mode
- **THEN** ScanResult.Vulnerabilities SHALL contain one entry per CVE with all fields populated
- **AND** ScanResult.Counts and Fixable SHALL still be populated as before

#### Scenario: Standalone mode uses cache-dir
- **WHEN** scanning an image in standalone mode (no server URL configured)
- **THEN** Trivy is invoked with `--cache-dir <dir>` and `--skip-java-db-update`

#### Scenario: Server mode scan with vulnerabilities
- **WHEN** scanning an image in server mode (server URL configured)
- **THEN** Trivy is invoked with `--server <url>` (no `--cache-dir`, no `--skip-java-db-update`)
- **AND** ScanResult.Vulnerabilities SHALL contain one entry per CVE with all fields populated

#### Scenario: Scan with no vulnerabilities
- **WHEN** scanning a clean image in either mode
- **THEN** ScanResult.Vulnerabilities SHALL be empty

#### Scenario: CLI constructor accepts server URL
- **WHEN** NewCLI is called with a non-empty serverURL parameter
- **THEN** the CLI operates in server mode
- **AND** cacheDir is ignored for scan operations
