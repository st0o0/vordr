## ADDED Requirements

### Requirement: Trivy DB age metric
The metrics collector SHALL expose `vordr_trivy_db_updated_timestamp_seconds` gauge with the Unix timestamp of the last Trivy DB update.

#### Scenario: DB age metric present
- **WHEN** /metrics is scraped after at least one DB check
- **THEN** `vordr_trivy_db_updated_timestamp_seconds` SHALL be present with a non-zero value

#### Scenario: DB check failure
- **WHEN** `trivy version --format json` fails
- **THEN** the metric SHALL retain its last known value and the failure SHALL be logged

### Requirement: Periodic DB version check
A background goroutine SHALL run `trivy version --format json` periodically (every 30 minutes) and extract the VulnerabilityDB.UpdatedAt timestamp.

#### Scenario: Periodic check
- **WHEN** 30 minutes elapse
- **THEN** the DB check SHALL execute and update the stored timestamp

#### Scenario: Startup check
- **WHEN** vordr starts
- **THEN** the DB check SHALL execute once immediately
