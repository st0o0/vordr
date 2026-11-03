## MODIFIED Requirements

### Requirement: Image state lifecycle
The store key SHALL remain `ImageRef.String()` (Image:Tag) for deduplication. The ImageRef stored SHALL include the ImageID field for scan targeting.

#### Scenario: Add image with ImageID
- **WHEN** an image is added with ImageID `sha256:52d09...`
- **THEN** the stored ImageRef SHALL contain the ImageID
- **AND** Snapshot() SHALL return the ImageID in each ImageState

#### Scenario: Images() returns ImageID
- **WHEN** the scheduler requests known images for rescan
- **THEN** the returned ImageRef values SHALL include ImageID
