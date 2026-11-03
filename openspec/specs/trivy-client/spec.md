## MODIFIED Requirements

### Requirement: Trivy scan execution
The Trivy CLI client SHALL parse the full vulnerability details from Trivy JSON output, not just severity counts. ScanResult SHALL include a Vulnerabilities slice with ID, PkgName, InstalledVersion, FixedVersion, and Severity per CVE.

#### Scenario: Successful scan with vulnerabilities
- **WHEN** scanning an image with vulnerabilities
- **THEN** ScanResult.Vulnerabilities SHALL contain one entry per CVE with all fields populated
- **AND** ScanResult.Counts and Fixable SHALL still be populated as before

#### Scenario: Scan with no vulnerabilities
- **WHEN** scanning a clean image
- **THEN** ScanResult.Vulnerabilities SHALL be empty
