## ADDED Requirements

### Requirement: Thread-safe image state management
The store SHALL manage image state with a sync.RWMutex protecting a map of image references to ImageState. It SHALL support concurrent reads (metrics scrape) and writes (worker updates) safely.

#### Scenario: Concurrent read and write
- **WHEN** a worker updates image state while the metrics handler reads a snapshot
- **THEN** neither operation SHALL block the other unsafely (RWMutex)

### Requirement: Image state lifecycle
The store SHALL track images through states: Pending, Scanning, Scanned, Error. State transitions SHALL be enforced.

#### Scenario: New image added
- **WHEN** an image is added to the store
- **THEN** its status SHALL be Pending with an empty container list

#### Scenario: Scan starts
- **WHEN** a worker begins scanning an image
- **THEN** the status SHALL transition from Pending or Scanned or Error to Scanning

#### Scenario: Scan succeeds
- **WHEN** a worker completes a scan with results
- **THEN** the status SHALL be Scanned with vulnerability counts, fixable counts, scan timestamp, and scan duration populated

#### Scenario: Scan fails
- **WHEN** a worker scan fails
- **THEN** the status SHALL be Error with the error message stored

### Requirement: Container tracking per image
The store SHALL maintain a list of container names for each image. Adding/removing containers SHALL NOT trigger rescans.

#### Scenario: Add container to existing image
- **WHEN** a new container starts with an image already in the store
- **THEN** the container name SHALL be appended to the image's container list

#### Scenario: Remove last container
- **WHEN** the last container for an image is removed
- **THEN** the image entry SHALL be deleted from the store

### Requirement: Snapshot for metrics
The store SHALL provide a Snapshot() method that returns a copy of all image states without holding the lock during serialization.

#### Scenario: Snapshot returns current state
- **WHEN** Snapshot() is called
- **THEN** it SHALL return a slice of ImageState copies representing the current state
- **AND** the read lock SHALL be released before returning

### Requirement: Known images listing
The store SHALL provide a method to return all currently tracked image references for the scheduler to queue rescans.

#### Scenario: Images for rescan
- **WHEN** the scheduler requests known images
- **THEN** the store SHALL return ImageRef values for all tracked images
