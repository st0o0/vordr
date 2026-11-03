## MODIFIED Requirements

### Requirement: Metrics HTTP server
The metrics server SHALL expose /details in addition to /metrics and /healthz. The Trivy DB age metric SHALL be included in the Prometheus collector output.

#### Scenario: Details endpoint registered
- **WHEN** the metrics server starts
- **THEN** GET /details SHALL be routable and return JSON

#### Scenario: DB age in metrics
- **WHEN** /metrics is scraped
- **THEN** vordr_trivy_db_updated_timestamp_seconds SHALL be present
