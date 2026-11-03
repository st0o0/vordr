## ADDED Requirements

### Requirement: Trivy server scan API
The Trivy client SHALL send HTTP POST requests to the Trivy server scan endpoint and parse the JSON response to extract vulnerability counts by severity.

#### Scenario: Successful scan
- **WHEN** scanning `grafana/grafana:13.3.0` and Trivy returns results
- **THEN** the client SHALL return vulnerability counts per severity (CRITICAL, HIGH, MEDIUM, LOW, UNKNOWN) and fixable counts per severity

#### Scenario: Image with no vulnerabilities
- **WHEN** scanning an image and Trivy returns empty results
- **THEN** the client SHALL return zero counts for all severities

#### Scenario: Trivy server unreachable
- **WHEN** the Trivy server is not reachable
- **THEN** the client SHALL return an error within the configured timeout

#### Scenario: Trivy returns error response
- **WHEN** the Trivy server returns a non-200 HTTP status
- **THEN** the client SHALL return an error with the status code and response body

### Requirement: HTTP client configuration
The Trivy client SHALL use the configured VORDR_TRIVY_ADDR and a reasonable request timeout.

#### Scenario: Custom Trivy address
- **WHEN** VORDR_TRIVY_ADDR is `http://trivy.local:8080`
- **THEN** the client SHALL send requests to `http://trivy.local:8080`

#### Scenario: Request timeout
- **WHEN** a scan request exceeds 5 minutes
- **THEN** the client SHALL cancel the request and return a timeout error
