## ADDED Requirements

### Requirement: Exponential backoff retry for failed scans
Failed scans SHALL be automatically retried with exponential backoff instead of waiting for the next scheduled rescan.

#### Scenario: First failure
- **WHEN** a scan fails
- **THEN** the image SHALL be requeued after 30 seconds

#### Scenario: Repeated failures
- **WHEN** a scan fails again after retry
- **THEN** the backoff SHALL double (30s, 1m, 2m, 4m, 8m, 15m max)

#### Scenario: Successful scan after retries
- **WHEN** a scan succeeds after previous failures
- **THEN** the backoff counter SHALL reset to zero

#### Scenario: Max backoff
- **WHEN** the backoff exceeds VORDR_RETRY_MAX_BACKOFF (default 15m)
- **THEN** it SHALL be capped at the max value

#### Scenario: Scheduled rescan during retry
- **WHEN** the scheduler queues a rescan for an image already in the retry queue
- **THEN** the retry SHALL proceed normally (no deduplication needed — SetScanning handles it)
