## ADDED Requirements

### Requirement: DB freshness monitoring
The DBChecker SHALL be disabled when vordr operates in server mode. In standalone mode, it SHALL continue to check the local DB and emit `vordr_trivy_db_updated_timestamp_seconds`.

#### Scenario: Standalone mode - DB check active
- **WHEN** vordr runs in standalone mode
- **THEN** DBChecker periodically runs `trivy version --format json` and emits the DB timestamp metric

#### Scenario: Server mode - DB check disabled
- **WHEN** vordr runs in server mode
- **THEN** DBChecker.Run() returns immediately without executing any checks
- **AND** `vordr_trivy_db_updated_timestamp_seconds` is not emitted
