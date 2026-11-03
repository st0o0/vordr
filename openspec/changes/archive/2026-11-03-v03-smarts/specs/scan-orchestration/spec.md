## MODIFIED Requirements

### Requirement: Periodic rescan
The scheduler SHALL skip excluded images when queuing rescans.

#### Scenario: Excluded image during rescan
- **WHEN** the scheduler iterates known images and one matches an exclude pattern
- **THEN** it SHALL NOT be sent to the scan channel

### Requirement: Worker pool scan with delta detection
The worker SHALL perform delta detection after each successful scan and log new CVEs. Failed scans SHALL be sent to the retry queue.

#### Scenario: Scan failure triggers retry
- **WHEN** a scan fails
- **THEN** the worker SHALL send the ImageRef to the retry queue
