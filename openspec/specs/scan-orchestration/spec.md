## MODIFIED Requirements

### Requirement: Structured logging
All components SHALL use structured slog logging with an `event` attribute for machine-parseable log events.

#### Scenario: Scan lifecycle events
- **WHEN** a scan completes
- **THEN** the log entry SHALL include `event=scan_completed`, `image`, `duration`, and `counts` attributes

#### Scenario: Discovery events
- **WHEN** a new image is discovered
- **THEN** the log entry SHALL include `event=image_discovered`, `image`, and `container` attributes

#### Scenario: Scheduled rescan
- **WHEN** the scheduler triggers a rescan
- **THEN** the log entry SHALL include `event=scheduled_rescan` and `images` count
