## MODIFIED Requirements

### Requirement: Image state lifecycle
The store SHALL store vulnerability details ([]Vulnerability) alongside counts in ImageState. UpdateResult SHALL accept vulnerabilities. Snapshot SHALL copy vulnerabilities.

#### Scenario: Store vulnerability details
- **WHEN** a worker updates scan results with vulnerability details
- **THEN** the ImageState SHALL contain the full vulnerability list

#### Scenario: Snapshot includes vulnerabilities
- **WHEN** Snapshot() is called
- **THEN** each ImageState SHALL include a copy of its vulnerabilities slice
