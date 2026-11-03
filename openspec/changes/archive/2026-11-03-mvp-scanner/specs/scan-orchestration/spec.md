## ADDED Requirements

### Requirement: Worker pool bounded parallelism
The worker pool SHALL run N goroutines (configured via VORDR_SCAN_WORKERS) that read from a shared scan channel. Each worker SHALL scan the image via the Trivy client, then update the store with results.

#### Scenario: Concurrent scanning
- **WHEN** VORDR_SCAN_WORKERS is 3 and 5 images are queued
- **THEN** at most 3 images SHALL be scanned concurrently

#### Scenario: Worker receives ImageRef
- **WHEN** a worker reads an ImageRef from the scan channel
- **THEN** it SHALL set the image status to Scanning in the store
- **AND** call the Trivy client to scan the image
- **AND** update the store with scan results (status Scanned) or error (status Error)

#### Scenario: Scan error
- **WHEN** the Trivy client returns an error for an image
- **THEN** the worker SHALL set the image status to Error with the error message
- **AND** continue processing the next ImageRef from the channel

#### Scenario: Skip already-scanning image
- **WHEN** a worker reads an ImageRef that is already in Scanning status
- **THEN** the worker SHALL skip it without calling Trivy

### Requirement: Staggered startup scan
The watcher SHALL send initial scan requests with a configurable delay (VORDR_SCAN_STAGGER) between each image to avoid overwhelming Trivy on startup.

#### Scenario: Staggered scheduling
- **WHEN** vordr starts with 10 discovered images and VORDR_SCAN_STAGGER is 30s
- **THEN** scan requests SHALL be sent approximately 30s apart
- **AND** the first image SHALL be sent immediately

#### Scenario: Shutdown during stagger
- **WHEN** the context is cancelled during staggered scheduling
- **THEN** remaining images SHALL NOT be sent and the watcher SHALL return

### Requirement: Periodic rescan
The scheduler SHALL periodically queue all known images for rescanning at VORDR_SCAN_INTERVAL (default 6h).

#### Scenario: Scheduled rescan
- **WHEN** VORDR_SCAN_INTERVAL elapses and the store contains 5 images
- **THEN** 5 ImageRef values SHALL be sent to the scan channel

#### Scenario: Scheduler graceful shutdown
- **WHEN** the context is cancelled
- **THEN** the scheduler SHALL stop the ticker and return nil

### Requirement: Errgroup lifecycle
All components (metrics server, watcher, scheduler, worker pool) SHALL run as goroutines in an errgroup. If any component fails, the shared context SHALL be cancelled and all components SHALL shut down.

#### Scenario: Component failure
- **WHEN** the watcher fails with an error
- **THEN** the errgroup context SHALL be cancelled
- **AND** all other components SHALL receive the cancellation and shut down

#### Scenario: Clean shutdown
- **WHEN** SIGTERM is received
- **THEN** the context SHALL be cancelled and errgroup.Wait SHALL return nil after all components return
