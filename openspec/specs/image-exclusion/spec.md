## ADDED Requirements

### Requirement: Image exclude patterns
The system SHALL support a comma-separated list of glob patterns (VORDR_EXCLUDE_IMAGES) to skip images from scanning. Matching uses Go's path.Match against "image:tag".

#### Scenario: Excluded image on startup
- **WHEN** vordr discovers a container with image matching an exclude pattern
- **THEN** the image SHALL NOT be added to the store or queued for scanning

#### Scenario: Excluded image on event
- **WHEN** a container start event matches an exclude pattern
- **THEN** the image SHALL be ignored

#### Scenario: Excluded image on scheduled rescan
- **WHEN** the scheduler queues images and one matches an exclude pattern
- **THEN** it SHALL be skipped

#### Scenario: No exclude patterns configured
- **WHEN** VORDR_EXCLUDE_IMAGES is empty
- **THEN** all images SHALL be scanned (no filtering)

#### Scenario: Multiple patterns
- **WHEN** VORDR_EXCLUDE_IMAGES is `aquasec/trivy:*,*devcontainer*`
- **THEN** images matching either pattern SHALL be excluded
