## MODIFIED Requirements

### Requirement: Previous vulnerability tracking
The store SHALL track the previous set of CVE IDs per image for delta detection.

#### Scenario: After UpdateResult
- **WHEN** UpdateResult is called with new vulnerabilities
- **THEN** the previous vulnerability IDs SHALL be saved before overwriting with new ones

#### Scenario: PreviousVulnIDs in Snapshot
- **WHEN** Snapshot() is called
- **THEN** each ImageState SHALL include PreviousVulnIDs for delta comparison by the worker
