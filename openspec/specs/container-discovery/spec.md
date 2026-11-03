## MODIFIED Requirements

### Requirement: Initial container listing
The watcher SHALL skip images matching exclude patterns during initial container listing.

#### Scenario: Excluded image at startup
- **WHEN** a running container's image matches an exclude pattern
- **THEN** it SHALL NOT be added to the store or queued for scanning
- **AND** a debug log with event=image_excluded SHALL be emitted
