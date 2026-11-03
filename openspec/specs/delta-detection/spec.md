## ADDED Requirements

### Requirement: CVE delta detection on rescan
The worker SHALL compare new CVE IDs against previously known CVE IDs after each scan. Only newly appeared CVEs SHALL be logged as a delta event.

#### Scenario: First scan (no previous data)
- **WHEN** an image is scanned for the first time
- **THEN** no delta event SHALL be emitted (all CVEs are new, no baseline to compare)

#### Scenario: Rescan with new CVEs
- **WHEN** a rescan finds CVEs not present in the previous scan
- **THEN** an `event=new_cves_detected` log entry SHALL be emitted with the count and list of new CVE IDs

#### Scenario: Rescan with no new CVEs
- **WHEN** a rescan finds the same or fewer CVEs
- **THEN** no delta event SHALL be emitted

### Requirement: Severity-based alert threshold
New CVEs SHALL be classified as "alert-worthy" if their severity meets or exceeds VORDR_ALERT_SEVERITY (default HIGH). The delta log SHALL include an `alert_worthy` count.

#### Scenario: New CRITICAL CVE with threshold HIGH
- **WHEN** a new CVE has severity CRITICAL and threshold is HIGH
- **THEN** it SHALL be counted as alert-worthy

#### Scenario: New LOW CVE with threshold HIGH
- **WHEN** a new CVE has severity LOW and threshold is HIGH
- **THEN** it SHALL NOT be counted as alert-worthy
