## MODIFIED Requirements

### Requirement: Metrics HTTP server
The metrics server SHALL expose /details in addition to /metrics and /healthz. The Trivy DB age metric SHALL be included in the Prometheus collector output. The server SHALL serve the embedded web dashboard on GET /.

#### Scenario: Details endpoint registered
- **WHEN** the metrics server starts
- **THEN** GET /details SHALL be routable and return JSON

#### Scenario: DB age in metrics
- **WHEN** /metrics is scraped
- **THEN** vordr_trivy_db_updated_timestamp_seconds SHALL be present

#### Scenario: Dashboard served on root
- **WHEN** GET / is requested
- **THEN** the server SHALL return the embedded HTML dashboard with Content-Type text/html
